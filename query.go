package chem

import (
	"errors"
	"fmt"
	"strings"
)

type BasicSelect struct {
	columns []string
	table   string
}

func (s BasicSelect) QueryString() (string, error) {
	if len(s.columns) == 0 {
		return "", errors.New("No columns specified")
	}

	if len(s.table) == 0 {
		return "", errors.New("No table specified")
	}

	return fmt.Sprintf(
		"SELECT %v FROM %v",
		strings.Join(s.columns, ", "),
		s.table,
	), nil
}

type BasicInsert struct {
	columns []string
	table   string
}

func thisMany(num int, thing string) (out []string) {
	for i := 0; i < num; i++ {
		out = append(out, thing)
	}
	return
}

func (i BasicInsert) QueryString() (string, error) {
	if len(i.columns) == 0 {
		return "", errors.New("No columns specified")
	}

	if len(i.table) == 0 {
		return "", errors.New("No table specified")
	}

	return fmt.Sprintf(
		"INSERT INTO %v (%v) VALUES (%v)",
		i.table,
		strings.Join(i.columns, ", "),
		strings.Join(thisMany(len(i.columns), "?"), ", "),
	), nil
}
