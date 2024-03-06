package repository

import "testing"

func TestEqual(t *testing.T) {
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
}

func TestGreaterThan(t *testing.T) {
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
}

func TestGreaterThanEqual(t *testing.T) {
	field := "id"
	value := 10
	c := GreaterThan(field, value)
	q := c.GetQuery()
	v := c.GetValues()
	eq := field + " >= ?"
	ev := []int{value}

	if eq == q && len(v) == 1 && v[0] == value {
		t.Logf("GreaterThan(\"%s\", \"%d\") PASS. Expected query: \"%s\", value: %d. Got \"%s\", %d", field, value, eq, ev, q, v)
	} else {
		t.Errorf("GreaterThan(\"%s\", \"%d\") FAILED. Expected query: \"%s\", value: %d. Got \"%s\", %d", field, value, eq, ev, q, v)
	}
}

func TestLessThan(t *testing.T) {
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
}

func TestLessThanEqual(t *testing.T) {
	field := "id"
	value := 10
	c := LessThan(field, value)
	q := c.GetQuery()
	v := c.GetValues()
	eq := field + " <= ?"
	ev := []int{value}

	if eq == q && len(v) == 1 && v[0] == value {
		t.Logf("LessThan(\"%s\", \"%d\") PASS. Expected query: \"%s\", value: %d. Got \"%s\", %d", field, value, eq, ev, q, v)
	} else {
		t.Errorf("LessThan(\"%s\", \"%d\") FAILED. Expected query: \"%s\", value: %d. Got \"%s\", %d", field, value, eq, ev, q, v)
	}
}
