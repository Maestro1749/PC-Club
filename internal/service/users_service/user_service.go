package users_service

import (
	"mvp/internal/models"
	"mvp/internal/repository/users_repository"
	"mvp/pkg"
	"time"

	"go.uber.org/zap"
)

type UserServise struct {
	repo   users_repository.UserRepository
	logger *zap.Logger
}

func NewService(repo users_repository.UserRepository, logger *zap.Logger) *UserServise {
	return &UserServise{
		repo:   repo,
		logger: logger,
	}
}

func (s *UserServise) RegisterUser(newUser models.NewUserDTO) error {
	if err := pkg.ValidateUsername(newUser.Username); err != nil {
		return err
	}

	if err := pkg.ValidateEmail(newUser.Email); err != nil {
		return err
	}

	birthday, err := time.Parse("02.01.2006", newUser.Birthday)
	if err != nil {
		return models.ErrInvalidBirtday
	}

	// Сейчас поддерживается только номера РФ. Позже добавить возможность зарегистрировать другие страны.
	if err := pkg.ValidatePhoneNumber(newUser.PhoneNumber, "RU"); err != nil {
		return err
	}

	if err := pkg.ValidatePassword(newUser.Password); err != nil {
		return err
	}

	hash, err := pkg.HashPassword(newUser.Password)
	if err != nil {
		return err
	}

	user := models.User{
		Username:    newUser.Username,
		Fullname:    newUser.Fullname,
		Email:       newUser.Email,
		PhoneNumber: newUser.PhoneNumber,
		Password:    string(hash),
		Birthday:    birthday,
		Registered:  time.Now(),
	}

	if err := s.repo.Create(&user); err != nil {
		return err
	}

	return nil
}

func (s *UserServise) LoginUser(userInfo models.LoginUserDTO) (*models.User, error) {
	if err := pkg.ValidateUsername(userInfo.Username); err != nil {
		return nil, err // invalid username
	}

	if err := pkg.ValidatePassword(userInfo.Password); err != nil {
		return nil, err // wrong password
	}

	hash, err := s.repo.GetHashByUsername(userInfo.Username)
	if err != nil {
		return nil, err // not found
	}

	if err := pkg.CheckPassword(userInfo.Password, hash); err != nil {
		return nil, err // wrong password
	}

	user, err := s.repo.LoginUser(userInfo.Username, string(hash))
	if err != nil {
		return nil, err
	}

	return user, nil
}
