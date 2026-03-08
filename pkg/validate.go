package pkg

import (
	"errors"
	"net/mail"
	"regexp"

	"github.com/nyaruka/phonenumbers"
)

func ValidateUsername(username string) error {
	re := regexp.MustCompile(`^[a-zA-Z0-9_]{4,20}$`)
	if !re.MatchString(username) {
		return errors.New("Invalid username.")
	}

	return nil
}

func ValidateEmail(email string) error {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return errors.New("Invalid email.")
	}

	return nil
}

func ValidatePhoneNumber(phoneNumber string, region string) error {
	num, err := phonenumbers.Parse(phoneNumber, region)
	if err != nil {
		return errors.New("Invalid phone number.")
	}

	if !phonenumbers.IsValidNumber(num) {
		return errors.New("Invalid phone number.")
	}

	return nil
}

func ValidatePassword(password string) error {
	if len(password) < 8 || len(password) > 100 {
		return errors.New("Password is too short.")
	}

	return nil
}
