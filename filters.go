package chem

import (
	"fmt"
	"strings"
)

type Filter interface {
	binds() []interface{}
	toBooleanExpression() string
}

const (
	equalsOperator = "=="
)

type ValueFilter struct {
	column   Column
	operator string
	value    interface{}
}

func (f ValueFilter) toBooleanExpression() string {
	return fmt.Sprintf("%v %v ?", f.column.toColumnExpression(), f.operator)
}

func (f ValueFilter) binds() []interface{} {
	return []interface{}{f.value}
}

type BooleanOperatorFilter struct {
	operator string
	filters  []Filter
}

func (f BooleanOperatorFilter) toBooleanExpression() string {
	expressionList := make([]string, len(f.filters))
	for i, filter := range f.filters {
		expressionList[i] = filter.toBooleanExpression()
	}
	expression := strings.Join(
		expressionList,
		fmt.Sprintf(" %v ", f.operator),
	)
	if expression == "" {
		return expression
	}
	// wrap the expression to make sure precedence is what user expects
	return fmt.Sprintf("( %v )", expression)
}

func (f BooleanOperatorFilter) binds() []interface{} {
	out := make([]interface{}, 0, len(f.filters))
	for _, filter := range f.filters {
		out = append(out, filter.binds()...)
	}
	return out
}

const (
	andOperator = "AND"
)

func And(filters ...Filter) Filter {
	return BooleanOperatorFilter{
		operator: andOperator,
		filters:  filters,
	}
}
