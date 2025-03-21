package usecase

import (
	"api_quiz/dto"
	"api_quiz/internal/repository"
	"api_quiz/utils/helper"
)

type AuthUseCase interface {
	Login(dto *dto.Login) (string, error)
	Register(input *dto.Register) error
	DeleteUser(id uint) error

	ValidateUser(id uint) error
}

type authUseCase struct {
	authRepo repository.AuthRepository
}

func NewAuthUseCase(authRepo repository.AuthRepository) AuthUseCase {
	return &authUseCase{authRepo}
}

func (u *authUseCase) Login(dto *dto.Login) (string, error) {

	if !helper.IsValidEmail(dto.Email) {
		return "", helper.ErrInvalidEmail
	}

	user, err := u.authRepo.Login(dto)
	if err != nil {
		return "", err
	}

	if !helper.ComparePassword(user.Password, dto.Password) {
		return "", err
	}

	token, err := helper.GenerateJWTLogin(user.ID, user.Email, user.IsVerified)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (u *authUseCase) Register(input *dto.Register) error {
	if !helper.IsValidEmail(input.Email) {
		return helper.ErrInvalidEmail
	}
	hashed, err := helper.HashPassword(input.Password)
	if err != nil {
		return err
	}

	input.Password = hashed
	user, err := u.authRepo.Register(input)
	if err != nil {
		return err
	}

	token, err := helper.GenerateJWTRegister(user.ID, user.Email)
	if err != nil {
		return err
	}

	helper.SendEmail(user.Email, token)
	return nil
}

func (u *authUseCase) DeleteUser(id uint) error {
	return u.authRepo.DeleteUser(id)
}

func (u *authUseCase) ValidateUser(id uint) error {
	return u.authRepo.ValidateUser(id)
}
