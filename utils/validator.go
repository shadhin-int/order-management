package utils

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"regexp"
	"time"
)

func ValidatePhoneNumber(phoneNumber string) bool {
	pattern := `^(01)[3-9][0-9]{8}$`
	match, _ := regexp.MatchString(pattern, phoneNumber)
	return match

}

func GetValidationErrors(err error) map[string][]string {
	errors := make(map[string][]string)
	if ve, ok := err.(validator.ValidationErrors); ok {
		for _, e := range ve {
			field := e.Field()
			errors[field] = append(errors[field], "The "+field+" field is required.")
		}
	}
	return errors
}

func GenerateConsignmentID() string {
	return fmt.Sprintf("DA%s%s",
		time.Now().Format("060102"),
		RandomString(6))
}

func RandomString(n int) string {
	var letter = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letter[time.Now().UnixNano()%int64(len(letter))]
	}
	return string(b)
}
