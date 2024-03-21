package db_repository

import (
	"fmt"
	"strings"
)

type join struct {
	conditions []Condition
	separator  string
}

func (s join) GetQuery() string {
	queries := make([]string, len(s.conditions))

	for _, spec := range s.conditions {
		queries = append(queries, spec.GetQuery())
	}

	return "(" + strings.Join(queries, fmt.Sprintf(" %s ", s.separator)) + ")"
}

func (s join) GetValues() []any {
	values := make([]any, len(s.conditions))

	for _, spec := range s.conditions {
		values = append(values, spec.GetValues()...)
	}

	return values
}

func And(conditions ...Condition) Condition {
	return join{
		conditions: conditions,
		separator:  "AND",
	}
}

func Or(conditions ...Condition) Condition {
	return join{
		conditions: conditions,
		separator:  "OR",
	}
}

type not struct {
	c Condition
}

func (s not) GetQuery() string {
	return fmt.Sprintf(" NOT (%s)", s.c.GetQuery())
}

func (s not) GetValues() []any {
	return s.c.GetValues()
}

func Not(condition Condition) Condition {
	return not{c: condition}
}

type binary[T any] struct {
	field    string
	operator string
	value    T
}

func (s binary[T]) GetQuery() string {
	return fmt.Sprintf("%s %s ?", s.field, s.operator)
}

func (s binary[T]) GetValues() []any {
	return []any{s.value}
}

func Equal[T any](field string, value T) Condition {
	return binary[T]{
		field:    field,
		operator: "=",
		value:    value,
	}
}

func NotEqual[T any](field string, value T) Condition {
	return binary[T]{
		field:    field,
		operator: "!=",
		value:    value,
	}
}

func GreaterThan[T comparable](field string, value T) Condition {
	return binary[T]{
		field:    field,
		operator: ">",
		value:    value,
	}
}

func GreaterOrEqual[T comparable](field string, value T) Condition {
	return binary[T]{
		field:    field,
		operator: ">=",
		value:    value,
	}
}

func LessThan[T comparable](field string, value T) Condition {
	return binary[T]{
		field:    field,
		operator: "<",
		value:    value,
	}
}

func LessOrEqual[T comparable](field string, value T) Condition {
	return binary[T]{
		field:    field,
		operator: "<=",
		value:    value,
	}
}

func Contains(field string, value string) Condition {
	return binary[string]{
		field:    field,
		operator: "like",
		value:    "%" + value + "%",
	}
}

type margin[T any] struct {
	field    string
	operator string
	from     T
	to       T
}

func (s margin[T]) GetQuery() string {
	return fmt.Sprintf("%s %s ? AND ?", s.field, s.operator)
}

func (s margin[T]) GetValues() []any {
	return []any{s.from, s.to}
}

func Between[T any](field string, from T, to T) Condition {
	return margin[T]{
		field:    field,
		operator: "BETWEEN",
		from:     from,
		to:       to,
	}
}

func NotBetween[T any](field string, from T, to T) Condition {
	return margin[T]{
		field:    field,
		operator: "NOT BETWEEN",
		from:     from,
		to:       to,
	}
}

type str string

func (s str) GetQuery() string {
	return string(s)
}

func (s str) GetValues() []any {
	return []any{}
}

func IsNull(field string) Condition {
	return str(fmt.Sprintf("%s IS NULL", field))
}

type array[T any] struct {
	field    string
	operator string
	values   []T
}

func (s array[T]) GetQuery() string {
	return fmt.Sprintf("%s IN (?)", s.field)
}

func (s array[T]) GetValues() []any {
	return []any{s.values}
}

func In[T string | uint](field string, values []T) Condition {
	return array[T]{
		field:    field,
		operator: "in",
		values:   values,
	}
}

func NotIn[T string | int](field string, values []T) Condition {
	return array[T]{
		field:    field,
		operator: "not in",
		values:   values,
	}
}
