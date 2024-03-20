package hashing

import (
	"errors"
	"sync"
)

type Contact interface {
	Make(str string) (string, error)
	Check(hashed, plain string) bool
}

const (
	DriverBcrypt = "bcrypt"
	DriverArgon2 = "argon2"
)

type Config struct {
	Driver string
	Bcrypt BcryptConfig
	Argon2 Argon2Config
}

var (
	hashingInstance Contact
	hashingOnce     sync.Once
)

func New(c Config) (Contact, error) {
	var err error
	hashingOnce.Do(func() {
		switch c.Driver {
		case DriverBcrypt:
			hashingInstance = newBcrypt(c.Bcrypt)
			break
		case DriverArgon2:
			hashingInstance = newArgon2(c.Argon2)
			break
		default:
			err = errors.New("hashing driver is invalid")
		}
	})
	return hashingInstance, err
}
