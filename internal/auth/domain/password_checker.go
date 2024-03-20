package domain

import "github.com/kurneo/go-template/pkg/hashing"

type PasswordChecker struct {
	s hashing.Contact
}

func (c PasswordChecker) Check(hashed, plain string) bool {
	return c.s.Check(hashed, plain)
}

func NewPasswordChecker(s hashing.Contact) *PasswordChecker {
	return &PasswordChecker{
		s: s,
	}
}
