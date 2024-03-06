package validator

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"reflect"
	"strings"
)

func ValidateStruct(stc interface{}) map[string][]string {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
		}
	}()
	errors := make(map[string][]string, 0)
	validate := validator.New()
	validate.RegisterTagNameFunc(func(structField reflect.StructField) string {
		name := strings.SplitN(structField.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
	err := validate.Struct(stc)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			rules := []string{err.Tag()}
			errors[strings.ToLower(err.Field())] = rules
		}
	}
	return errors
}

func ValidateValue(value interface{}, rules string) []string {
	validate := validator.New()

	err := validate.Var(value, rules)

	if err != nil {
		err := err.(validator.ValidationErrors)[0]
		rules := []string{err.Tag()}
		return rules
	}

	return nil
}
