package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

func HashPassword(password string) string {
	sum := sha256.Sum256([]byte(password))
	return hex.EncodeToString(sum[:])
}

func CheckPassword(password, hash string) bool {
	return HashPassword(password) == hash
}
