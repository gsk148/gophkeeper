package services

import (
	"errors"

	"github.com/gsk148/gophkeeper/internal/app/server/storage"
	"github.com/gsk148/gophkeeper/internal/pkg/enc"
)

type UserService struct {
	db storage.IUserRepository
}

func NewUserService(db storage.IUserRepository) UserService {
	return UserService{db: db}
}

func (s UserService) AddUser(req AuthReq) error {
	user := getUserFromRequest(req)
	userExist, err := doesUserExist(s.db, user)
	if err != nil {
		return err
	}
	if userExist {
		return errors.New("user with the specified name already exists")
	}

	hash, err := enc.HashPassword(user.Password)
	if err != nil {
		return nil
	}

	_, err = s.db.AddUser(user.Name, hash)
	return err
}

func (s UserService) GetUser(user AuthReq) (storage.User, error) {
	su, err := s.db.GetUserByName(user.Name)
	if err != nil {
		return storage.User{}, err
	}

	if !enc.VerifyPassword(user.Password, su.Password) {
		return storage.User{}, errors.New("user not found")
	}
	return su, nil
}

func getUserFromRequest(r AuthReq) storage.User {
	return storage.User{
		Name:     r.Name,
		Password: r.Password,
	}
}

func doesUserExist(db storage.IUserRepository, user storage.User) (bool, error) {
	su, err := db.GetUserByName(user.Name)
	if err != nil && err.Error() != "user not found" {
		return false, err
	}
	return su.ID != "", nil
}
