package chem

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
)

type Table interface {
	Name() string
	Type() reflect.Type
}

type InsertStmt struct {
	table Table
}

func Insert(table Table) InsertStmt {
	return InsertStmt{table: table}
}

func binds(num int) []string {
	out := make([]string, num)
	for i := 0; i < num; i++ {
		out[i] = "?"
	}
	return out
}

func (i InsertStmt) Values(tx *sql.Tx, value interface{}) (sql.Result, error) {
	reflection := reflect.ValueOf(value)
	reflectedType := reflection.Type()
	if tableType := i.table.Type(); reflectedType != tableType {
		err := IncorrectTypeError{
			Got:      reflectedType,
			Expected: tableType,
		}
		return BadResult{err}, err
	}
	columns := make([]string, reflectedType.NumField())
	values := make([]interface{}, reflectedType.NumField())
	for i := 0; i < reflectedType.NumField(); i++ {
		structField := reflectedType.Field(i)
		columns[i] = func(structField reflect.StructField) string {
			name := structField.Tag.Get("chem")
			if name == "" {
				name = structField.Name
			}
			return name
		}(structField)

		values[i] = reflection.Field(i).Interface()
	}

	queryString := fmt.Sprintf(
		"INSERT INTO %v (%v) VALUES (%v)",
		i.table.Name(),
		strings.Join(columns, ", "),
		strings.Join(binds(len(columns)), ", "),
	)

	return tx.Exec(queryString, values...)
}
