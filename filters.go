package chem

import (
	"fmt"
	"strings"
)

type Filter interface {
	binds() []interface{}
	toBooleanExpression() string
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

func And(filters ...Filter) Filter {
	return BooleanOperatorFilter{
		operator: "AND",
		filters:  filters,
	}
}
