package utils

import (
	"regexp"
	"strings"

	"github.com/sethvargo/go-password/password"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
func ComparePasswords(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
func GenerateRandomPassword() string {
	strRes, err := password.Generate(64, 10, 10, false, false)
	if err != nil {
		return GenerateRandomSecureString(64)
	}
	return strRes
}
func IsStrongPassword(password string) bool {
	if len(password) < 8 {
		return false
	}

	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasDigit := regexp.MustCompile(`[0-9]`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`).MatchString(password)

	if !hasUpper || !hasLower || !hasDigit || !hasSpecial {
		return false
	}

	weakPasswords := []string{"password", "123456", "qwerty", "12345678", "12345", "1234"}
	for _, weak := range weakPasswords {
		if strings.Contains(strings.ToLower(password), weak) {
			return false
		}
	}

	return true
}
