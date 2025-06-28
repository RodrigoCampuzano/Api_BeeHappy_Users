package auth

import (
	"crypto/sha256"
	"encoding/hex"
)

func HashPasswordWithSecret(password, secret string) string {
	hash := sha256.Sum256([]byte(password + secret))
	return hex.EncodeToString(hash[:])
}
