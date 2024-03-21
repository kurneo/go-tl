package db_repository

import (
	"reflect"
	"testing"
)

func TestBinary(t *testing.T) {
	t.Run("TestEqual", func(t *testing.T) {
		field := "name"
		value := "test"
		c := Equal(field, value)
		q := c.GetQuery()
		v := c.GetValues()
		eq := field + " = ?"
		ev := []string{value}

		if eq == q && len(v) == 1 && v[0] == value {
			t.Logf("Equal(\"%s\", \"%s\") PASS. Expected query: \"%s\", value: %s. Got \"%s\", %s", field, value, eq, ev, q, v)
		} else {
			t.Errorf("Equal(\"%s\", \"%s\") FAILED. Expected query: \"%s\", value: %s. Got \"%s\", %s", field, value, eq, ev, q, v)
		}
	})

	t.Run("TestNotEqual", func(t *testing.T) {
		field := "name"
		value := "test"
		c := NotEqual(field, value)
		q := c.GetQuery()
		v := c.GetValues()
		eq := field + " != ?"
		ev := []string{value}

		if eq == q && len(v) == 1 && v[0] == value {
			t.Logf("Equal(\"%s\", \"%s\") PASS. Expected query: \"%s\", value: %s. Got \"%s\", %s", field, value, eq, ev, q, v)
		} else {
			t.Errorf("Equal(\"%s\", \"%s\") FAILED. Expected query: \"%s\", value: %s. Got \"%s\", %s", field, value, eq, ev, q, v)
		}
	})

	t.Run("TestGreaterThan", func(t *testing.T) {
		field := "id"
		value := 10
		c := GreaterThan(field, value)
		q := c.GetQuery()
		v := c.GetValues()
		eq := field + " > ?"
		ev := []int{value}

		if eq == q && len(v) == 1 && v[0] == value {
			t.Logf("GreaterThan(\"%s\", \"%d\") PASS. Expected query: \"%s\", value: %d. Got \"%s\", %d", field, value, eq, ev, q, v)
		} else {
			t.Errorf("GreaterThan(\"%s\", \"%d\") FAILED. Expected query: \"%s\", value: %d. Got \"%s\", %d", field, value, eq, ev, q, v)
		}
	})

	t.Run("TestGreaterOrEqual", func(t *testing.T) {
		field := "id"
		value := 10
		c := GreaterOrEqual(field, value)
		q := c.GetQuery()
		v := c.GetValues()
		eq := field + " >= ?"
		ev := []int{value}

		if eq == q && len(v) == 1 && v[0] == value {
			t.Logf("GreaterOrEqual(\"%s\", \"%d\") PASS. Expected query: \"%s\", value: %d. Got \"%s\", %d", field, value, eq, ev, q, v)
		} else {
			t.Errorf("GreaterOrEqual(\"%s\", \"%d\") FAILED. Expected query: \"%s\", value: %d. Got \"%s\", %d", field, value, eq, ev, q, v)
		}
	})

	t.Run("TestLessThan", func(t *testing.T) {
		field := "id"
		value := 10
		c := LessThan(field, value)
		q := c.GetQuery()
		v := c.GetValues()
		eq := field + " < ?"
		ev := []int{value}

		if eq == q && len(v) == 1 && v[0] == value {
			t.Logf("LessThan(\"%s\", \"%d\") PASS. Expected query: \"%s\", value: %d. Got \"%s\", %d", field, value, eq, ev, q, v)
		} else {
			t.Errorf("LessThan(\"%s\", \"%d\") FAILED. Expected query: \"%s\", value: %d. Got \"%s\", %d", field, value, eq, ev, q, v)
		}
	})

	t.Run("TestLessOrEqual", func(t *testing.T) {
		field := "id"
		value := 10
		c := LessOrEqual(field, value)
		q := c.GetQuery()
		v := c.GetValues()
		eq := field + " <= ?"
		ev := []int{value}

		if eq == q && len(v) == 1 && v[0] == value {
			t.Logf("LessOrEqual(\"%s\", \"%d\") PASS. Expected query: \"%s\", value: %d. Got \"%s\", %d", field, value, eq, ev, q, v)
		} else {
			t.Errorf("LessOrEqual(\"%s\", \"%d\") FAILED. Expected query: \"%s\", value: %d. Got \"%s\", %d", field, value, eq, ev, q, v)
		}
	})
}

func TestJoin(t *testing.T) {
	t.Run("TestAnd", func(t *testing.T) {
		c1 := Equal("name", "test")
		c2 := NotEqual("age", 30)
		c := And(c1, c2)
		q := c.GetQuery()
		v := c.GetValues()
		eq := "(" + c1.GetQuery() + " AND " + c2.GetQuery() + ")"
		ev := append(append([]any{}, c1.GetValues()...), c2.GetValues()...)

		if eq == q && reflect.DeepEqual(v, ev) {
			t.Logf("And(c1, c2) PASS. Expected query: \"%s\", value: %v. Got \"%s\", %v", eq, ev, q, v)
		} else {
			t.Errorf("And(c1, c2) FAILED. Expected query: \"%s\", value: %v. Got \"%s\", %v", eq, ev, q, v)
		}
	})

	t.Run("TestOr", func(t *testing.T) {
		c1 := Equal("name", "test")
		c2 := NotEqual("age", 30)
		c := Or(c1, c2)
		q := c.GetQuery()
		v := c.GetValues()
		eq := "(" + c1.GetQuery() + " OR " + c2.GetQuery() + ")"
		ev := append(append([]any{}, c1.GetValues()...), c2.GetValues()...)

		if eq == q && reflect.DeepEqual(v, ev) {
			t.Logf("Or(c1, c2) PASS. Expected query: \"%s\", value: %v. Got \"%s\", %v", eq, ev, q, v)
		} else {
			t.Errorf("Or(c1, c2) FAILED. Expected query: \"%s\", value: %v. Got \"%s\", %v", eq, ev, q, v)
		}
	})
}

func TestMargin(t *testing.T) {
	t.Run("TestBetween", func(t *testing.T) {
		field := "age"
		from := 15
		to := 60
		c := Between(field, from, to)
		q := c.GetQuery()
		v := c.GetValues()
		eq := field + " BETWEEN ? AND ?"
		ev := []any{from, to}

		if eq == q && reflect.DeepEqual(v, ev) {
			t.Logf("Between(\"%s\", %d, %d) PASS. Expected query: \"%s\", value: %v. Got \"%s\", %v", field, from, to, eq, ev, q, v)
		} else {
			t.Errorf("Between(\"%s\", %d, %d) FAILED. Expected query: \"%s\", value: %v. Got \"%s\", %v", field, from, to, eq, ev, q, v)
		}
	})

	t.Run("TestNotBetween", func(t *testing.T) {
		field := "age"
		from := 15
		to := 60
		c := NotBetween(field, from, to)
		q := c.GetQuery()
		v := c.GetValues()
		eq := field + " NOT BETWEEN ? AND ?"
		ev := []any{from, to}

		if eq == q && reflect.DeepEqual(v, ev) {
			t.Logf("NotBetween(\"%s\", %d, %d) PASS. Expected query: \"%s\", value: %v. Got \"%s\", %v", field, from, to, eq, ev, q, v)
		} else {
			t.Errorf("NotBetween(\"%s\", %d, %d) FAILED. Expected query: \"%s\", value: %v. Got \"%s\", %v", field, from, to, eq, ev, q, v)
		}
	})
}
