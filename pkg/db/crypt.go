package db

import (
	"crypto/sha256"
	"fmt"
	"os"
)

func encrypt(s string) string {
	hash := sha256.Sum256([]byte(os.Getenv("SALT") + s))
	return fmt.Sprintf("%x", hash)
}

func compare(s string, givenHash string) bool {
	hash := sha256.Sum256([]byte(os.Getenv("SALT") + s))
	return givenHash == fmt.Sprintf("%x", hash)
}
