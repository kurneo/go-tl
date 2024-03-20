package hashing

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

type Bcrypt struct {
	c BcryptConfig
}

type BcryptConfig struct {
	Cost int
}

func (b Bcrypt) Make(str string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(str), b.c.Cost)
	return string(bytes), err
}

func (b Bcrypt) Check(hashed, plain string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain)) == nil
}

func newBcrypt(c BcryptConfig) *Bcrypt {
	if c.Cost == 0 {
		c.Cost = 14
	}
	return &Bcrypt{
		c: c,
	}
}

type Argon2 struct {
	c Argon2Config
}

type Argon2Config struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	SaltLength  uint32
	KeyLength   uint32
}

func (a Argon2) Make(str string) (string, error) {
	return a.generateFromPassword(str)
}

func (a Argon2) Check(hashed, plain string) bool {
	check, err := a.comparePasswordAndHash(plain, hashed)
	if !check || err != nil {
		return false
	}
	return true
}

func (a Argon2) generateFromPassword(password string) (string, error) {
	// Generate a cryptographically secure random salt.
	salt, err := a.generateRandomBytes(a.c.SaltLength)
	if err != nil {
		return "", err
	}
	// Pass the plaintext password, salt and parameters to the argon2.IDKey
	// function. This will generate a hash of the password using the Argon2id
	// variant.
	hash := argon2.IDKey([]byte(password), salt, a.c.Iterations, a.c.Memory, a.c.Parallelism, a.c.KeyLength)

	// Base64 encode the salt and hashed password.
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	return fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, a.c.Memory, a.c.Iterations, a.c.Parallelism, b64Salt, b64Hash), nil
}

func (a Argon2) generateRandomBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (a Argon2) comparePasswordAndHash(password, encodedHash string) (match bool, err error) {
	// Extract the parameters, salt and derived key from the encoded password
	// hash.
	p, salt, hash, err := a.decodeHash(encodedHash)
	if err != nil {
		return false, err
	}

	// Derive the key from the other password using the same parameters.
	otherHash := argon2.IDKey([]byte(password), salt, p.Iterations, p.Memory, p.Parallelism, p.KeyLength)

	// Check that the contents of the hashed passwords are identical. Note
	// that we are using the subtle.ConstantTimeCompare() function for this
	// to help prevent timing attacks.
	if subtle.ConstantTimeCompare(hash, otherHash) == 1 {
		return true, nil
	}
	return false, nil
}

func (a Argon2) decodeHash(encodedHash string) (*Argon2Config, []byte, []byte, error) {
	vals := strings.Split(encodedHash, "$")
	if len(vals) != 6 {
		return nil, nil, nil, errors.New("the encoded hash is not in the correct format")
	}

	var version int
	_, err := fmt.Sscanf(vals[2], "v=%d", &version)
	if err != nil {
		return nil, nil, nil, err
	}
	if version != argon2.Version {
		return nil, nil, nil, errors.New("incompatible version of argon2")
	}

	p := &Argon2Config{}
	_, err = fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &p.Memory, &p.Iterations, &p.Parallelism)
	if err != nil {
		return nil, nil, nil, err
	}

	salt, err := base64.RawStdEncoding.Strict().DecodeString(vals[4])
	if err != nil {
		return nil, nil, nil, err
	}
	p.SaltLength = uint32(len(salt))

	hash, err := base64.RawStdEncoding.Strict().DecodeString(vals[5])
	if err != nil {
		return nil, nil, nil, err
	}
	p.KeyLength = uint32(len(hash))

	return p, salt, hash, nil
}

func newArgon2(c Argon2Config) *Argon2 {
	if c.Memory == 0 {
		c.Memory = 64 * 1024
	}

	if c.Iterations == 0 {
		c.Iterations = 3
	}

	if c.Parallelism == 0 {
		c.Parallelism = 2
	}

	if c.SaltLength == 0 {
		c.SaltLength = 16
	}

	if c.KeyLength == 0 {
		c.KeyLength = 32
	}

	return &Argon2{
		c: c,
	}
}
