package chem

import (
	"fmt"
	"reflect"
)

type IncorrectTypeError struct {
	Got      reflect.Type
	Expected reflect.Type
}

func (e IncorrectTypeError) Error() string {
	return fmt.Sprintf("unexpected type: got %v, expected %v", e.Got, e.Expected)
}

type BadResult struct {
	Err error
}

func (b BadResult) Error() string {
	return b.Err.Error()
}

func (b BadResult) LastInsertId() (int64, error) {
	return -1, b.Err
}

func (b BadResult) RowsAffected() (int64, error) {
	return -1, b.Err
}

type NonSliceError struct {
	Type reflect.Type
}

func (e NonSliceError) Error() string {
	return fmt.Sprintf("expecting a slice type, got %v", e.Type)
}
