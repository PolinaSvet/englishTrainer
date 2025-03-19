package pass

import (
	"errors"
	"regexp"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword хеширует пароль
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash проверяет, соответствует ли пароль хешу
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// ValidatePassword проверяет сложность пароля
func ValidatePassword(password string) error {
	var errorsList []string

	// Минимальная длина пароля
	if len(password) < 8 {
		errorsList = append(errorsList, "password must be at least 8 characters long")
	}

	// Проверка на наличие цифр
	if match, _ := regexp.MatchString(`[0-9]`, password); !match {
		errorsList = append(errorsList, "password must contain at least one digit")
	}

	// Проверка на наличие заглавных букв
	if match, _ := regexp.MatchString(`[A-Z]`, password); !match {
		errorsList = append(errorsList, "password must contain at least one uppercase letter")
	}

	// Проверка на наличие строчных букв
	if match, _ := regexp.MatchString(`[a-z]`, password); !match {
		errorsList = append(errorsList, "password must contain at least one lowercase letter")
	}

	// Проверка на наличие специальных символов
	if match, _ := regexp.MatchString(`[!@#$%^&*()_+{}:"<>?]`, password); !match {
		errorsList = append(errorsList, "password must contain at least one special character")
	}

	if len(errorsList) > 0 {
		return errors.New(strings.Join(errorsList, "; "))
	}

	return nil
}
