package chem

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
)

type SelectStmt struct {
	columns []Column
	filters []Filter
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

func flattenValues(values []interface{}) (out []interface{}) {
	for _, value := range values {
		reflection := reflect.ValueOf(value).Elem()
		reflectType := reflection.Type()
		switch reflectType.Kind() {
		case reflect.Struct:
			for i := 0; i < reflection.NumField(); i++ {
				out = append(out, reflection.Field(i).Addr().Interface())
			}
		default:
			out = append(out, value)
		}
	}
	return
}

func makeWhereClause(f Filter) string {
	expression := f.toBooleanExpression()
	if expression == "" {
		return ""
	}
	return fmt.Sprintf("WHERE %v", expression)
}

func (s SelectStmt) constructSQL() string {
	return strings.Join(
		filterStringSlice(
			fmt.Sprintf(
				"SELECT %v FROM %v",
				strings.Join(toColumnExpressions(s.columns), ", "),
				strings.Join(toTableNames(s.columns), ", "),
			),
			makeWhereClause(AND(s.filters...)),
		),
		" ",
	)
}

func (s SelectStmt) One(tx *sql.Tx, values ...interface{}) error {
	return tx.QueryRow(
		s.constructSQL(),
		AND(s.filters...).binds()...,
	).Scan(flattenValues(values)...)
}
