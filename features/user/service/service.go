package service

import (
	"errors"
	"library_api/features/user"
	"library_api/helper/enkrip"
	"strings"
)

type UserService struct {
	repo user.Repository
	hash enkrip.HashInterface
}

func New(r user.Repository, h enkrip.HashInterface) user.Service {
	return &UserService{
		repo: r,
		hash: h,
	}
}

// Login implements user.Service.
func (us *UserService) Login(phone string, password string) (user.User, error) {
	if phone == "" || password == "" {
		return user.User{}, errors.New("phone and password are required")
	}
	result, err := us.repo.Login(phone)

	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return user.User{}, errors.New("data tidak ditemukan")
		}
		return user.User{}, errors.New("data tidak ditemukan")
	}

	err = us.hash.Compare(result.Password, password)

	if err != nil {
		return user.User{}, errors.New("password salah")
	}

	return result, nil
}

// Register implements user.Service.
func (us *UserService) Register(newUser user.User) (user.User, error) {
	if newUser.Name == "" {
		return user.User{}, errors.New("name cannot be empty")
	}
	if newUser.Phone == "" {
		return user.User{}, errors.New("phone cannot be empty")
	}
	if newUser.Password == "" {
		return user.User{}, errors.New("password cannot be empty")
	}
	// enkripsi password
	ePassword, err := us.hash.HashPassword(newUser.Password)

	if err != nil {
		return user.User{}, errors.New("terdapat masalah saat memproses enkripsi password")
	}

	newUser.Password = ePassword
	result, err := us.repo.Register(newUser)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return user.User{}, errors.New("data telah terdaftar pada sistem")
		}
		return user.User{}, errors.New("terjadi kesalahan pada sistem")
	}

	return result, nil
}