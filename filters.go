package chem

import (
	"fmt"
)

type IntegerFilter struct {
	col   Column
	value int
}

func (f IntegerFilter) toBooleanExpression() string {
	return fmt.Sprintf("%v == ?", f.col.toColumnExpression())
}

func (f IntegerFilter) binds() (out []interface{}) {
	out = append(out, f.value)
	return
}

type Filter interface {
	binds() []interface{}
	toBooleanExpression() string
}
