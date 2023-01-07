package service

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/hex"
)

const saltSize = 16

func generateRandomSalt(saltSize int) ([]byte, error) {
	var salt = make([]byte, saltSize)

	_, err := rand.Read(salt)

	if err != nil {
		return nil, err
	}

	return salt, nil
}

func hashPassword(password string, salt []byte) string {
	pwdBytes := []byte(password)

	sha512Hasher := sha512.New()

	pwdBytes = append(pwdBytes, salt...)

	sha512Hasher.Write(pwdBytes)

	hashedPwdBytes := sha512Hasher.Sum(nil)

	return hex.EncodeToString(hashedPwdBytes)
}

func doPasswordsMatch(hashedPassword, currPassword string, salt []byte) bool {
	var currPasswordHash = hashPassword(currPassword, salt)

	return hashedPassword == currPasswordHash
}
