package services

import (
	"back/jwt"
	"back/repository"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Signup(user repository.User) error
	Login(email, password string) (*repository.User, string, error)
	GetUserByID(id string) (*repository.User, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(r repository.UserRepository) UserService {
	return &userService{repo: r}
}

func (u *userService) Signup(user repository.User) error {
	// if user already availble
	existing, _ := u.repo.FindByEmail(user.Email)
	if existing != nil {
		return errors.New("user already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	return u.repo.Create(user)
}

func (u *userService) Login(email, password string) (*repository.User, string, error) {
	user, err := u.repo.FindByEmail(email)
	if err != nil {
		return nil, "", errors.New("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, "", errors.New("invalid password")
	}
	token, err := jwt.GenerateJWT(user.ID)
	if err != nil {
		return nil, "", errors.New("failed to generate token")
	}

	return user, token, nil
}

func (u *userService) GetUserByID(id string) (*repository.User, error) {
	return u.repo.FindByID(id)
}
