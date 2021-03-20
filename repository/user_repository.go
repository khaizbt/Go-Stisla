package repository

import (
	"errors"
	"goshop/config"
	"goshop/model"

	"gorm.io/gorm"
)

type (
	UserRepository interface {
		FindUserByEmail(email string) (model.User, error)
		FindByID(ID int) (model.User, error)
		UpdateProfile(user model.User) (model.User, error)
		CreateUser(user model.User) (model.User, error)
		SaveOtp(otp model.UserCode) error
		CheckOtp(email, otp string) (bool, error)
		UpdatePassword(email, password string) error
		ChangeStatusUser(email string) error
	}

	repository struct {
		db *gorm.DB
	}
)

func NewUserRepository() *repository {
	return &repository{config.GetDB()}
}

func (r *repository) FindUserByEmail(email string) (model.User, error) {
	var user model.User
	err := r.db.Where("email = ?", email).First(&user).Error

	if err != nil {
		return user, err
	}

	return user, nil
}

func (r *repository) FindByID(ID int) (model.User, error) {
	var user model.User

	err := r.db.Where("id = ?", ID).First(&user).Error

	if err != nil {
		return user, err
	}

	return user, nil
}

func (r *repository) UpdateProfile(user model.User) (model.User, error) {
	err := r.db.Save(&user).Error

	if err != nil {
		return user, err
	}

	return user, nil
}

func (r *repository) CreateUser(user model.User) (model.User, error) {
	trx := r.db.Begin()
	err := trx.Create(&user).Error

	if err != nil {
		trx.Rollback()
		return user, err
	}
	trx.Commit()
	return user, nil
}

func (r *repository) SaveOtp(otp model.UserCode) error {
	err := r.db.Create(&otp).Error

	if err != nil {
		return err
	}

	return nil
}

func (r *repository) CheckOtp(email, otp string) (bool, error) {
	var otpCode model.UserCode
	err := r.db.Where("email = ?", email).Where("code = ?", otp).Last(&otpCode).Error

	if err != nil {
		return false, err
	}

	if otpCode.ID == 0 {
		return false, errors.New("Wrong Code!")
	}

	return true, nil
}

func (r *repository) UpdatePassword(email, password string) error {

	err := r.db.Model(&model.User{}).Where("email = ?", email).Update("password", password).Error

	if err != nil {
		return err
	}

	return nil
}

func (r *repository) ChangeStatusUser(email string) error {
	err := r.db.Model(&model.User{}).Where("email = ?", email).Update("status", "ACTIVE").Error

	if err != nil {
		return err
	}

	return nil
}
