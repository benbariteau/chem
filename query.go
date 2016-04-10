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
