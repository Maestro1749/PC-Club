package pkg

import (
	"mvp/internal/models"
	"net/mail"
	"regexp"

	"github.com/nyaruka/phonenumbers"
	"golang.org/x/crypto/bcrypt"
)

func ValidateUsername(username string) error {
	re := regexp.MustCompile(`^[a-zA-Z0-9_]{4,20}$`)
	if !re.MatchString(username) {
		return models.ErrInvalidUsername
	}

	return nil
}

func ValidateEmail(email string) error {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return models.ErrInvalidEmail
	}

	return nil
}

func ValidatePhoneNumber(phoneNumber string, region string) error {
	num, err := phonenumbers.Parse(phoneNumber, region)
	if err != nil {
		return models.ErrInvalidPhoneNumber
	}

	if !phonenumbers.IsValidNumber(num) {
		return models.ErrInvalidPhoneNumber
	}

	return nil
}

func ValidatePassword(password string) error {
	if len(password) < 8 {
		return models.ErrWrongPassword
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
		return models.ErrWrongPassword
	}

	return nil
}
