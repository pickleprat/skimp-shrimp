package _util

import (
	"encoding/base64"
	"errors"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func ConvertStringToUint(input string) (uint, error) {
	if input == "" {
		return 0, errors.New("input string is empty")
	}
	convertedUint, err := strconv.ParseUint(input, 10, 0)
	if err != nil {
		return 0, err
	}
	return uint(convertedUint), nil
}

func BytesToBase64String(input []byte) string {
	encodedString := base64.StdEncoding.EncodeToString(input)
	return encodedString
}

func URLBuilder(basePath string, params ...string) string {
	if len(params)%2 != 0 {
		panic("Invalid number of parameters. Must be even.")
	}

	if len(params) == 0 {
		return basePath
	}

	var sb strings.Builder
	sb.WriteString(basePath)
	sb.WriteString("?")

	for i := 0; i < len(params); i += 2 {
		if i > 0 {
			sb.WriteString("&")
		}
		sb.WriteString(params[i])
		sb.WriteString("=")
		sb.WriteString(params[i+1])
	}
	return sb.String()
}

func IsValidPhoneNumber(phone string) bool {
	pattern := `^\d{3}-\d{3}-\d{4}$`
	regex := regexp.MustCompile(pattern)
	return regex.MatchString(phone)
}

func IsValidEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	regex := regexp.MustCompile(pattern)
	return regex.MatchString(email)
}

func GenerateRandomToken(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	token := make([]byte, length)
	rand.Seed(time.Now().UnixNano())
	for i := range token {
		token[i] = charset[rand.Intn(len(charset))]
	}
	return string(token)
}

func StringWithDefault(input string, defaultValue string) string {
	if input == "" {
		return defaultValue
	}
	return input
}

func ConditionalString(condition bool, option1 string, option2 string) string {
	if condition {
		return option1
	}
	return option2
}
