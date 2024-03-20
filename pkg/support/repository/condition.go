package repository

import (
	"fmt"
	"strings"
)

type join struct {
	conditions []Condition
	separator  string
}

func (s join) GetQuery() string {
	queries := make([]string, 0, len(s.conditions))

	for _, spec := range s.conditions {
		queries = append(queries, spec.GetQuery())
	}

	return "(" + strings.Join(queries, fmt.Sprintf(" %s ", s.separator)) + ")"
}

func (s join) GetValues() []any {
	values := make([]any, 0)

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
	Condition
}

func (s not) GetQuery() string {
	return fmt.Sprintf(" NOT (%s)", s.Condition.GetQuery())
}

func Not(condition Condition) Condition {
	return not{condition}
}

type binaryOperator[T any] struct {
	field    string
	operator string
	value    T
}

func (s binaryOperator[T]) GetQuery() string {
	return fmt.Sprintf("%s %s ?", s.field, s.operator)
}

func (s binaryOperator[T]) GetValues() []any {
	return []any{s.value}
}

func Equal[T any](field string, value T) Condition {
	return binaryOperator[T]{
		field:    field,
		operator: "=",
		value:    value,
	}
}

func GreaterThan[T comparable](field string, value T) Condition {
	return binaryOperator[T]{
		field:    field,
		operator: ">",
		value:    value,
	}
}

func GreaterOrEqual[T comparable](field string, value T) Condition {
	return binaryOperator[T]{
		field:    field,
		operator: ">=",
		value:    value,
	}
}

func LessThan[T comparable](field string, value T) Condition {
	return binaryOperator[T]{
		field:    field,
		operator: "<",
		value:    value,
	}
}

func LessOrEqual[T comparable](field string, value T) Condition {
	return binaryOperator[T]{
		field:    field,
		operator: ">=",
		value:    value,
	}
}

func Contains(field string, value string) Condition {
	op := "like"
	return binaryOperator[string]{
		field:    field,
		operator: op,
		value:    "%" + value + "%",
	}
}

type fromToOperator[T any] struct {
	field    string
	operator string
	from     T
	to       T
}

func (s fromToOperator[T]) GetQuery() string {
	return fmt.Sprintf("%s %s ? AND ?", s.field, s.operator)
}

func (s fromToOperator[T]) GetValues() []any {
	return []any{s.from, s.to}
}

func Between[T any](field string, from T, to T) Condition {
	return fromToOperator[T]{
		field:    field,
		operator: "between",
		from:     from,
		to:       to,
	}
}

func NotBetween[T any](field string, from T, to T) Condition {
	return fromToOperator[T]{
		field:    field,
		operator: "not between",
		from:     from,
		to:       to,
	}
}

type str string

func (s str) GetQuery() string {
	return string(s)
}

func (s str) GetValues() []any {
	return nil
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
