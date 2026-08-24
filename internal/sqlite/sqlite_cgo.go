//go:build cgo

package sqlite

/*
#cgo LDFLAGS: -lsqlite3
#include <sqlite3.h>
#include <stdlib.h>

static int bind_text_transient(sqlite3_stmt* stmt, int index, const char* value, int size) {
    return sqlite3_bind_text(stmt, index, value, size, SQLITE_TRANSIENT);
}

static int bind_blob_transient(sqlite3_stmt* stmt, int index, const void* value, int size) {
    return sqlite3_bind_blob(stmt, index, value, size, SQLITE_TRANSIENT);
}
*/
import "C"

import (
	"context"
	"errors"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sync"
	"unsafe"
)

type ValueType int

const (
	Null ValueType = iota
	Integer
	Float
	Text
	Blob
)

type Value struct {
	Type    ValueType
	Integer int64
	Float   float64
	Text    string
	Blob    []byte
}

func (v Value) StringValue() (string, bool) {
	if v.Type != Text {
		return "", false
	}
	return v.Text, true
}

func (v Value) Int64Value() (int64, bool) {
	if v.Type != Integer {
		return 0, false
	}
	return v.Integer, true
}

type Row map[string]Value

type DB struct {
	mu     sync.Mutex
	handle *C.sqlite3
	closed bool
}

type Conn struct {
	handle *C.sqlite3
}

func Open(path string) (*DB, error) {
	if path == "" {
		return nil, errors.New("sqlite path is required")
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("resolve sqlite path: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(abs), 0o755); err != nil {
		return nil, fmt.Errorf("create sqlite directory: %w", err)
	}
	cpath := C.CString(abs)
	defer C.free(unsafe.Pointer(cpath))
	var handle *C.sqlite3
	flags := C.SQLITE_OPEN_READWRITE | C.SQLITE_OPEN_CREATE | C.SQLITE_OPEN_FULLMUTEX
	if rc := C.sqlite3_open_v2(cpath, &handle, C.int(flags), nil); rc != C.SQLITE_OK {
		message := "unknown sqlite open error"
		if handle != nil {
			message = C.GoString(C.sqlite3_errmsg(handle))
			_ = C.sqlite3_close_v2(handle)
		}
		return nil, fmt.Errorf("open sqlite database: %s (code %d)", message, int(rc))
	}
	C.sqlite3_extended_result_codes(handle, 1)
	db := &DB{handle: handle}
	if err := db.ExecScript(context.Background(), `
PRAGMA journal_mode=WAL;
PRAGMA foreign_keys=ON;
PRAGMA busy_timeout=5000;
PRAGMA synchronous=NORMAL;
`); err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}

func (db *DB) Close() error {
	db.mu.Lock()
	defer db.mu.Unlock()
	if db.closed {
		return nil
	}
	if rc := C.sqlite3_close_v2(db.handle); rc != C.SQLITE_OK {
		return sqliteError(db.handle, rc, "close database")
	}
	db.closed = true
	db.handle = nil
	return nil
}

func (db *DB) ExecScript(ctx context.Context, script string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	db.mu.Lock()
	defer db.mu.Unlock()
	if db.closed {
		return errors.New("sqlite database is closed")
	}
	cscript := C.CString(script)
	defer C.free(unsafe.Pointer(cscript))
	var errMsg *C.char
	rc := C.sqlite3_exec(db.handle, cscript, nil, nil, &errMsg)
	if rc != C.SQLITE_OK {
		message := C.GoString(errMsg)
		C.sqlite3_free(unsafe.Pointer(errMsg))
		return fmt.Errorf("execute sqlite script: %s (code %d)", message, int(rc))
	}
	return nil
}

func (db *DB) Exec(ctx context.Context, query string, args ...any) (int64, error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	if db.closed {
		return 0, errors.New("sqlite database is closed")
	}
	return (&Conn{handle: db.handle}).Exec(ctx, query, args...)
}

func (db *DB) Query(ctx context.Context, query string, args ...any) ([]Row, error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	if db.closed {
		return nil, errors.New("sqlite database is closed")
	}
	return (&Conn{handle: db.handle}).Query(ctx, query, args...)
}

func (db *DB) InTx(ctx context.Context, fn func(*Conn) error) (err error) {
	if fn == nil {
		return errors.New("transaction callback is nil")
	}
	db.mu.Lock()
	defer db.mu.Unlock()
	if db.closed {
		return errors.New("sqlite database is closed")
	}
	conn := &Conn{handle: db.handle}
	if _, err := conn.Exec(ctx, "BEGIN IMMEDIATE"); err != nil {
		return err
	}
	committed := false
	defer func() {
		if !committed {
			_, rollbackErr := conn.Exec(context.Background(), "ROLLBACK")
			if err == nil && rollbackErr != nil {
				err = rollbackErr
			}
		}
	}()
	if err := fn(conn); err != nil {
		return err
	}
	if _, err := conn.Exec(ctx, "COMMIT"); err != nil {
		return err
	}
	committed = true
	return nil
}

func (c *Conn) Exec(ctx context.Context, query string, args ...any) (int64, error) {
	if err := ctx.Err(); err != nil {
		return 0, err
	}
	stmt, err := prepare(c.handle, query)
	if err != nil {
		return 0, err
	}
	defer C.sqlite3_finalize(stmt)
	if err := bind(stmt, args); err != nil {
		return 0, err
	}
	rc := C.sqlite3_step(stmt)
	if rc != C.SQLITE_DONE && rc != C.SQLITE_ROW {
		return 0, sqliteError(c.handle, rc, "execute statement")
	}
	for rc == C.SQLITE_ROW {
		if err := ctx.Err(); err != nil {
			C.sqlite3_interrupt(c.handle)
			return 0, err
		}
		rc = C.sqlite3_step(stmt)
	}
	if rc != C.SQLITE_DONE {
		return 0, sqliteError(c.handle, rc, "finish statement")
	}
	return int64(C.sqlite3_changes(c.handle)), nil
}

func (c *Conn) Query(ctx context.Context, query string, args ...any) ([]Row, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	stmt, err := prepare(c.handle, query)
	if err != nil {
		return nil, err
	}
	defer C.sqlite3_finalize(stmt)
	if err := bind(stmt, args); err != nil {
		return nil, err
	}
	columnCount := int(C.sqlite3_column_count(stmt))
	columns := make([]string, columnCount)
	for i := range columns {
		columns[i] = C.GoString(C.sqlite3_column_name(stmt, C.int(i)))
	}
	var rows []Row
	for {
		if err := ctx.Err(); err != nil {
			C.sqlite3_interrupt(c.handle)
			return nil, err
		}
		rc := C.sqlite3_step(stmt)
		switch rc {
		case C.SQLITE_DONE:
			return rows, nil
		case C.SQLITE_ROW:
			row := make(Row, columnCount)
			for i, name := range columns {
				row[name] = columnValue(stmt, i)
			}
			rows = append(rows, row)
		default:
			return nil, sqliteError(c.handle, rc, "query statement")
		}
	}
}

func prepare(handle *C.sqlite3, query string) (*C.sqlite3_stmt, error) {
	cquery := C.CString(query)
	defer C.free(unsafe.Pointer(cquery))
	var stmt *C.sqlite3_stmt
	if rc := C.sqlite3_prepare_v2(handle, cquery, -1, &stmt, nil); rc != C.SQLITE_OK {
		return nil, sqliteError(handle, rc, "prepare statement")
	}
	if stmt == nil {
		return nil, errors.New("sqlite statement is empty")
	}
	return stmt, nil
}

func bind(stmt *C.sqlite3_stmt, args []any) error {
	want := int(C.sqlite3_bind_parameter_count(stmt))
	if want != len(args) {
		return fmt.Errorf("sqlite bind count mismatch: statement wants %d, got %d", want, len(args))
	}
	for i, arg := range args {
		index := C.int(i + 1)
		var rc C.int
		switch value := arg.(type) {
		case nil:
			rc = C.sqlite3_bind_null(stmt, index)
		case string:
			cvalue := C.CString(value)
			rc = C.bind_text_transient(stmt, index, cvalue, C.int(len(value)))
			C.free(unsafe.Pointer(cvalue))
		case []byte:
			if len(value) == 0 {
				rc = C.bind_blob_transient(stmt, index, nil, 0)
			} else {
				rc = C.bind_blob_transient(stmt, index, unsafe.Pointer(&value[0]), C.int(len(value)))
			}
		case int:
			rc = C.sqlite3_bind_int64(stmt, index, C.sqlite3_int64(value))
		case int64:
			rc = C.sqlite3_bind_int64(stmt, index, C.sqlite3_int64(value))
		case uint64:
			if value > math.MaxInt64 {
				return fmt.Errorf("sqlite cannot bind uint64 %d", value)
			}
			rc = C.sqlite3_bind_int64(stmt, index, C.sqlite3_int64(value))
		case bool:
			if value {
				rc = C.sqlite3_bind_int(stmt, index, 1)
			} else {
				rc = C.sqlite3_bind_int(stmt, index, 0)
			}
		case float64:
			rc = C.sqlite3_bind_double(stmt, index, C.double(value))
		default:
			return fmt.Errorf("unsupported sqlite bind type %T", arg)
		}
		if rc != C.SQLITE_OK {
			return fmt.Errorf("bind sqlite argument %d: code %d", i+1, int(rc))
		}
	}
	return nil
}

func columnValue(stmt *C.sqlite3_stmt, index int) Value {
	idx := C.int(index)
	switch C.sqlite3_column_type(stmt, idx) {
	case C.SQLITE_INTEGER:
		return Value{Type: Integer, Integer: int64(C.sqlite3_column_int64(stmt, idx))}
	case C.SQLITE_FLOAT:
		return Value{Type: Float, Float: float64(C.sqlite3_column_double(stmt, idx))}
	case C.SQLITE_TEXT:
		ptr := C.sqlite3_column_text(stmt, idx)
		size := C.sqlite3_column_bytes(stmt, idx)
		return Value{Type: Text, Text: C.GoStringN((*C.char)(unsafe.Pointer(ptr)), size)}
	case C.SQLITE_BLOB:
		ptr := C.sqlite3_column_blob(stmt, idx)
		size := C.sqlite3_column_bytes(stmt, idx)
		if size == 0 {
			return Value{Type: Blob, Blob: []byte{}}
		}
		return Value{Type: Blob, Blob: C.GoBytes(ptr, size)}
	default:
		return Value{Type: Null}
	}
}

func sqliteError(handle *C.sqlite3, rc C.int, operation string) error {
	return fmt.Errorf("%s: %s (code %d)", operation, C.GoString(C.sqlite3_errmsg(handle)), int(rc))
}
