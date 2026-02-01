package utils

import (
	"regexp"
	"strings"
)

// ValidateEmail - проверить формат email
func ValidateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// ValidatePassword - проверить силу пароля
func ValidatePassword(password string) (bool, string) {
	if len(password) < 6 {
		return false, "Password must be at least 6 characters long"
	}
	return true, ""
}

// ValidatePhone - проверить формат телефона (казахстанский +7...)
func ValidatePhone(phone string) bool {
	// Убрать все пробелы и дефисы
	phone = strings.ReplaceAll(phone, " ", "")
	phone = strings.ReplaceAll(phone, "-", "")

	// Проверить формат +7XXXXXXXXXX
	phoneRegex := regexp.MustCompile(`^\+7\d{10}$`)
	return phoneRegex.MatchString(phone)
}

// SanitizeString - очистить строку от лишних пробелов
func SanitizeString(s string) string {
	return strings.TrimSpace(s)
}
