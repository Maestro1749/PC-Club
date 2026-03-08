package pkg

import (
	"errors"
	"mvp/internal/models"
	"net/mail"
	"regexp"

	"github.com/nyaruka/phonenumbers"
	"golang.org/x/crypto/bcrypt"
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
	if len(password) < 8 {
		return models.ErrIncorrectPassword
	}

	return nil
}

func HashPassword(password string) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, models.ErrInternalServer
	}

	return hash, nil
}

func CheckPassword(password string, hash []byte) error {
	if err := bcrypt.CompareHashAndPassword(hash, []byte(password)); err != nil {
		return models.ErrIncorrectPassword
	}

	return nil
}
