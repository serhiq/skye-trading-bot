package chat

import (
	"regexp"
	"strconv"
)

func ValidateRussianPhoneNumber(phone string) (bool, string) {
	reg := regexp.MustCompile(`[\(\)\-\s]+`)
	phone = reg.ReplaceAllString(phone, "")

	if _, err := strconv.Atoi(phone); err != nil {
		return false, "номер содержит недопустимые символы"
	}
	return true, ""
}
