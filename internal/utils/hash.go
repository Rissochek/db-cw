package utils

import (
	"errors"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

func GenerateHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		zap.S().Errorf("failed to generate hash: %v", err)
		return "", errors.New("failed to generate hash")
	}
	return string(hash), nil
}