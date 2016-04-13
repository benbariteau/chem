package chem

import (
	"database/sql"
	"fmt"
	"strings"
)

type SelectStmt struct {
	columns []Column
	filters []Filter
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
	return IntegerFilter{
		col:   c,
		value: i,
	}
}

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

type Columnser interface {
	Columns() []Column
}

type Filter interface {
	binds() []interface{}
	toBooleanExpression() string
}

func Select(columnThings ...Columnser) SelectStmt {
	columns := make([]Column, 0, len(columnThings))
	for _, c := range columnThings {
		columns = append(columns, c.Columns()...)
	}

	return SelectStmt{columns: columns}
}

func (s SelectStmt) Where(filters ...Filter) SelectStmt {
	s.filters = filters
	return s
}

func toTableNames(columns []Column) []string {
	names := make(map[string]bool)
	for _, c := range columns {
		names[c.Table().Name()] = true
	}
	nameList := make([]string, 0, len(names))
	for name := range names {
		nameList = append(nameList, name)
	}
	return nameList
}

func toColumnExpressions(columns []Column) (out []string) {
	for _, c := range columns {
		out = append(out, c.toColumnExpression())
	}
	return
}

func toBooleanExpressions(filters []Filter) (out []string) {
	for _, filter := range filters {
		out = append(out, filter.toBooleanExpression())
	}
	return
}

func getAllBinds(filters []Filter) []interface{} {
	out := make([]interface{}, 0, len(filters))
	for _, f := range filters {
		out = append(out, f.binds()...)
	}
	return out
}

func (s SelectStmt) One(tx *sql.Tx, values ...interface{}) error {
	stmt := fmt.Sprintf(
		"SELECT %v FROM %v WHERE %v",
		strings.Join(toColumnExpressions(s.columns), ", "),
		strings.Join(toTableNames(s.columns), ", "),
		strings.Join(toBooleanExpressions(s.filters), " AND "),
	)

	return tx.QueryRow(stmt, getAllBinds(s.filters)...).Scan(values...)
}
