package chem

import (
	"database/sql"
	"reflect"
)

type Table interface {
	Name() string
	Type() reflect.Type
}

type Columnser interface {
	Columns() []Column
}

// DB is simply the common inteface between *sql.Tx and *sql.DB
type DB interface {
	Exec(query string, arg ...interface{}) (sql.Result, error)
	Prepare(query string) (*sql.Stmt, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}
