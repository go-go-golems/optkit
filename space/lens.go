package space

import "fmt"

type Lens[C, V any] struct {
	Get func(C) V
	Put func(C, V) (C, error)
}

func (l Lens[C, V]) Validate() error {
	if l.Get == nil {
		return fmt.Errorf("lens getter is nil")
	}
	if l.Put == nil {
		return fmt.Errorf("lens setter is nil")
	}
	return nil
}

func CheckLensLaws[C comparable, V comparable](lens Lens[C, V], base C, first, second V) error {
	if err := lens.Validate(); err != nil {
		return err
	}
	current := lens.Get(base)
	getPut, err := lens.Put(base, current)
	if err != nil {
		return fmt.Errorf("get-put failed: %w", err)
	}
	if getPut != base {
		return fmt.Errorf("get-put law violated")
	}
	putFirst, err := lens.Put(base, first)
	if err != nil {
		return fmt.Errorf("put first: %w", err)
	}
	if got := lens.Get(putFirst); got != first {
		return fmt.Errorf("put-get law violated: got %v, want %v", got, first)
	}
	putSecond, err := lens.Put(putFirst, second)
	if err != nil {
		return fmt.Errorf("put second: %w", err)
	}
	directSecond, err := lens.Put(base, second)
	if err != nil {
		return fmt.Errorf("direct put second: %w", err)
	}
	if putSecond != directSecond {
		return fmt.Errorf("put-put law violated")
	}
	return nil
}
