package domain

import (
	"github.com/kurneo/go-template/pkg/support/crypto"
)

type PasswordChecker struct {
}

func (c PasswordChecker) Check(hashed, plain string) bool {
	return crypto.Check(hashed, plain) == nil
}

func NewPasswordChecker() *PasswordChecker {
	return &PasswordChecker{}
}
