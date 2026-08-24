//go:build !cgo

package sqlite

import (
	"context"
	"errors"
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
type DB struct{}
type Conn struct{}

var errCGO = errors.New("SQLite support requires CGO")

func Open(string) (*DB, error)                                       { return nil, errCGO }
func (db *DB) Close() error                                          { return nil }
func (db *DB) ExecScript(context.Context, string) error              { return errCGO }
func (db *DB) Exec(context.Context, string, ...any) (int64, error)   { return 0, errCGO }
func (db *DB) Query(context.Context, string, ...any) ([]Row, error)  { return nil, errCGO }
func (db *DB) InTx(context.Context, func(*Conn) error) error         { return errCGO }
func (c *Conn) Exec(context.Context, string, ...any) (int64, error)  { return 0, errCGO }
func (c *Conn) Query(context.Context, string, ...any) ([]Row, error) { return nil, errCGO }
