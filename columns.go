package chem

import (
	"fmt"
)

type Columnser interface {
	Columns() []Column
}

type Column interface {
	Table() Table
	toColumnExpression() string
}

type BaseColumn struct {
	table Table
	name  string
}

func (c BaseColumn) toColumnExpression() string {
	return fmt.Sprintf("%v.%v", c.table.Name(), c.name)
}

func (c BaseColumn) Table() Table {
	return c.table
}

type IntegerColumn struct {
	BaseColumn
}

func (c IntegerColumn) Equals(i int) Filter {
	return ValueFilter{
		column:   c,
		operator: equalsOperator,
		value:    i,
	}
}
