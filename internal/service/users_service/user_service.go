package users_service

import (
	"mvp/internal/models"
	"mvp/internal/repository/users_repository"
	"mvp/pkg"
)

type UserServise struct {
	repo users_repository.UserRepository
}

func NewService(repo users_repository.UserRepository) *UserServise {
	return &UserServise{repo: repo}
}

func (s *UserServise) RegisterUser(newUser models.NewUserDTO) error {
	if err := pkg.ValidateUsername(newUser.Username); err != nil {
		return err
	}

	if err := pkg.ValidateEmail(newUser.Email); err != nil {
		return err
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
		Birthday:    newUser.Birthday,
	}

	if err := s.repo.Create(&user); err != nil {
		return err
	}

	return nil
}

func (s *UserServise) LoginUser(userInfo models.LoginUserDTO) (*models.User, error) {
	if err := pkg.ValidateUsername(userInfo.Username); err != nil {
		return nil, err
	}

	if err := pkg.ValidatePassword(userInfo.Password); err != nil {
		return nil, err
	}

	hash, err := s.repo.GetHashByUsername(userInfo.Username)
	if err != nil {
		return nil, err
	}

	if err := pkg.CheckPassword(userInfo.Password, hash); err != nil {
		return nil, err
	}

	user, err := s.repo.LoginUser(userInfo.Username, string(hash))
	if err != nil {
		return nil, err
	}

	return user, nil
}
