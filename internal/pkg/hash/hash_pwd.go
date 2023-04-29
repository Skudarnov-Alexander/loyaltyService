package hash

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/hex"
)

const saltSize = 16

func GenerateRandomSalt() ([]byte, error) {
	var salt = make([]byte, saltSize)

	if _, err := rand.Read(salt); err != nil {
		return nil, err
	}

	return salt, nil
}

type PassHasher struct {
	salt []byte
}

func New(salt []byte) *PassHasher {
	return &PassHasher{
		salt: salt,
	}
}

func (ph *PassHasher) Hash(password string) string {
	pwdBytes := []byte(password)
	pwdBytes = append(pwdBytes, ph.salt...)

	sha512Hasher := sha512.New()
	sha512Hasher.Write(pwdBytes)

	hashedPwdBytes := sha512Hasher.Sum(nil)

	return hex.EncodeToString(hashedPwdBytes)
}

func (ph *PassHasher) IsPwdsMatched(savedHashedPwd, pwd string, salt []byte) bool {
	return savedHashedPwd == ph.Hash(pwd)
}
