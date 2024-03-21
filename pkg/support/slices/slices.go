package slices

import "reflect"

func Map[T any, R any](arr []T, f func(v T) R) []R {
	var result []R
	for _, item := range arr {
		result = append(result, f(item))
	}
	return result
}

func Reduce[T any, R any](arr []T, f func(R, T) R) R {
	var result R
	for _, v := range arr {
		result = f(result, v)
	}
	return result
}

func RemoveAt[T any](slice []T, s int) []T {
	return append(slice[:s], slice[s+1:]...)
}

func Filter[T any](arr []T, f func(e T) bool) []T {
	result := make([]T, 0)
	for _, v := range arr {
		if f(v) {
			result = append(result, v)
		}
	}
	return result
}

func FirstBy[T any](arr []T, f func(item T) bool) *T {
	for _, item := range arr {
		if f(item) {
			return &item
		}
	}
	return nil
}

func First[T any](arr []T) *T {
	for _, item := range arr {
		return &item
	}
	return nil
}

func Unique[T any](arr []T) []T {
	keys := make(map[interface{}]bool)
	var list []T

	// If the key(values of the slice) is not equal
	// to the already present value in new slice (list)
	// then we append it. else we jump on another element.
	for _, entry := range arr {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func Contains[T comparable](s []T, e T) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func Compare[T any](a1, a2 []T) bool {
	return reflect.DeepEqual(a1, a2)
}

func Some[T any](arr []T, f func(item T) bool) bool {
	for _, v := range arr {
		if f(v) {
			return true
		}
	}
	return false
}

func Always[T any](arr []T, f func(item T) bool) bool {
	c := 0
	for _, v := range arr {
		if f(v) {
			c++
		}
	}
	return c == len(arr)
}
