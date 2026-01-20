package utils

import (
	"errors"
	"regexp"
)

func ValidatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 6 characters long")
	}

	var (
		upper   = regexp.MustCompile("[A-Z]")
		lower   = regexp.MustCompile("[a-z]")
		number  = regexp.MustCompile("[0-9]")
		special = regexp.MustCompile("[!@#%$^&*.]")
	)

	if !upper.MatchString(password) {
		return errors.New("password must contain at least one uppercase letter")
	}

	if !lower.MatchString(password) {
		return errors.New("password must contain at least one lowercase letter")
	}

	if !number.MatchString(password) {
		return errors.New("password must contain at least one number")
	}

	if !special.MatchString(password) {
		return errors.New("password must contain at least one special character")
	}

	return nil
}
