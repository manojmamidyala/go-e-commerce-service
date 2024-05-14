package utils

import (
	"mami/e-commerce/commons/logger"

	"golang.org/x/crypto/bcrypt"
)

func HashAndSalt(pass []byte) string {
	hashed, err := bcrypt.GenerateFromPassword(pass, bcrypt.DefaultCost)
	if err != nil {
		logger.Errorf("Failed to generate password: %v", err)
		return ""
	}

	return string(hashed)
}
