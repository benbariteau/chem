package chem

import (
	"fmt"
	"sort"
)

type Columnser interface {
	Columns() []Column
}

type Column interface {
	Table() Table
	toColumnExpression(withTableName bool) string
}

type columns struct {
	columns *[]Column
	less    func(left, right Column) bool
}

func (cols columns) Len() int {
	return len(*cols.columns)
}

func (cols columns) Swap(i, j int) {
	(*cols.columns)[i], (*cols.columns)[j] = (*cols.columns)[j], (*cols.columns)[i]
}

func (cols columns) Less(i, j int) bool {
	return cols.less((*cols.columns)[i], (*cols.columns)[j])
}

func lessColumnsByUnqualifiedName(left, right Column) bool {
	return left.toColumnExpression(false) < right.toColumnExpression(false)
}

func sortColumns(cols []Column) []Column {
	copiedSlice := cols[:]
	sorter := columns{
		columns: &copiedSlice,
		less:    lessColumnsByUnqualifiedName,
	}
	sort.Sort(sorter)
	return *sorter.columns
}

type BaseColumn struct {
	table Table
	name  string
}

func (c BaseColumn) toColumnExpression(withTableName bool) string {
	if withTableName {
		return fmt.Sprintf("%v.%v", c.table.Name(), c.name)
	}
	return c.name
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

type StringColumn struct {
	BaseColumn
}

func (c StringColumn) Equals(s string) Filter {
	return ValueFilter{
		column:   c,
		operator: equalsOperator,
		value:    s,
	}
}
