package chem

import (
	"database/sql"
	"fmt"
	"strings"
)

type UpdateStmt struct {
	table   Table
	filters []Filter
}

func Update(table Table) UpdateStmt {
	return UpdateStmt{
		table: table,
	}
}

func (stmt UpdateStmt) Where(filters ...Filter) UpdateStmt {
	stmt.filters = append(stmt.filters, filters...)
	return stmt
}

func toColumnsAndBinds(values map[Column]interface{}) (columns []Column, binds []interface{}) {
	columns = make([]Column, 0, len(values))
	for column := range values {
		columns = append(columns, column)
	}
	// sort so we can have predictable output
	columns = sortColumns(columns)

	for _, column := range columns {
		binds = append(binds, values[column])
	}
	return
}

func toSetColumnExpressions(columns []Column) []string {
	out := make([]string, 0, len(columns))
	for _, column := range columns {
		expression := fmt.Sprintf("%v = ?", column.toColumnExpression(false))
		out = append(out, expression)
	}
	return out
}

func allBinds(bindLists ...[]interface{}) []interface{} {
	out := make([]interface{}, 0)
	for _, list := range bindLists {
		for _, bind := range list {
			out = append(out, bind)
		}
	}
	return out
}

func (stmt UpdateStmt) Set(tx *sql.Tx, values map[Column]interface{}) (sql.Result, error) {
	if len(values) == 0 {
		return BadResult{ErrTooFewValues}, ErrTooFewValues
	}
	columns, bindValues := toColumnsAndBinds(values)
	combinedFilter := AND(stmt.filters...)
	queryString := strings.Join(
		filterStringSlice(
			fmt.Sprintf("UPDATE %v", stmt.table.Name()),
			fmt.Sprintf("SET %v", strings.Join(toSetColumnExpressions(columns), ", ")),
			makeWhereClause(combinedFilter, false),
		),
		" ",
	)
	preparedStmt, err := tx.Prepare(queryString)
	if err != nil {
		return BadResult{err}, err
	}
	return preparedStmt.Exec(
		allBinds(
			bindValues,
			combinedFilter.binds(),
		)...,
	)
}
