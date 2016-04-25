package chem

import (
	"fmt"
)

type Ordering interface {
	toOrderingExpression(fullyQualifyColumns bool) string
}

type ColumnOrdering struct {
	column     Column
	descending bool
}

func toOrderingKeyword(descending bool) string {
	if descending {
		return "DESC"
	}
	return "ASC"
}

func (o ColumnOrdering) toOrderingExpression(fullyQualifyColumns bool) string {
	return fmt.Sprintf(
		"%v %v",
		o.column.toColumnExpression(fullyQualifyColumns),
		toOrderingKeyword(o.descending),
	)
}
