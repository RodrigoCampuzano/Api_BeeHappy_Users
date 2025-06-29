package auth

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image/png"
	"time"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

type TwoFactorService struct{}

func NewTwoFactorService() *TwoFactorService {
	return &TwoFactorService{}
}

// GenerateSecret genera un nuevo secreto para 2FA
func (tfs *TwoFactorService) GenerateSecret(accountName, issuer string) (*otp.Key, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      issuer,
		AccountName: accountName,
		SecretSize:  32,
	})
	if err != nil {
		return nil, fmt.Errorf("error generando secreto 2FA: %w", err)
	}
	return key, nil
}

// ValidateToken valida un token TOTP
func (tfs *TwoFactorService) ValidateToken(token, secret string) bool {
	return totp.Validate(token, secret)
}

// GenerateQRCode genera un código QR como string base64
func (tfs *TwoFactorService) GenerateQRCode(key *otp.Key) (string, error) {
	var buf bytes.Buffer
	img, err := key.Image(256, 256)
	if err != nil {
		return "", fmt.Errorf("error generando imagen QR: %w", err)
	}

	err = png.Encode(&buf, img)
	if err != nil {
		return "", fmt.Errorf("error codificando imagen QR: %w", err)
	}

	return fmt.Sprintf("data:image/png;base64,%s", base64.StdEncoding.EncodeToString(buf.Bytes())), nil
}

// GenerateBackupCodes genera códigos de respaldo para 2FA
func (tfs *TwoFactorService) GenerateBackupCodes() []string {
	codes := make([]string, 10)
	for i := range codes {
		codes[i] = generateRandomCode(8)
	}
	return codes
}

// generateRandomCode genera un código aleatorio de la longitud especificada
func generateRandomCode(length int) string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	code := make([]byte, length)
	for i := range code {
		code[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(code)
}