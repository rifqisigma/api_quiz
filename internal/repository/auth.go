package repository

import (
	"api_quiz/dto"
	"api_quiz/entity"
	"api_quiz/utils/helper"

	"gorm.io/gorm"
)

type AuthRepository interface {
	Register(dto *dto.Register) (*entity.User, error)
	Login(input *dto.Login) (*entity.User, error)
	DeleteUser(id uint) error

	ValidateUser(id uint) error
}

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authRepository{db}
}

func (r *authRepository) Register(dto *dto.Register) (*entity.User, error) {
	user := entity.User{
		Email:    dto.Email,
		Password: dto.Password,
		Username: dto.Username,
	}
	if err := r.db.Create(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *authRepository) Login(input *dto.Login) (*entity.User, error) {
	var user entity.User
	err := r.db.Where("email = ?", input.Email).First(&user).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, helper.ErrUserNotFound
		}
	}

	return &user, nil
}

func (r *authRepository) DeleteUser(id uint) error {
	err := r.db.Delete(&entity.User{}, id)
	if err.RowsAffected == 0 {
		return helper.ErrUserNotFound
	}
	if err.Error != nil {
		return err.Error
	}
	return nil
}

func (r *authRepository) ValidateUser(id uint) error {
	err := r.db.Model(&entity.User{}).Where("id", id).Update("is_verified", true).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return helper.ErrUserNotFound
		}
		if err == gorm.ErrRecordNotFound {
			return err
		}
	}

	return nil
}
