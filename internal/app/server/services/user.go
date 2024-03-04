package services

import (
	"context"
	"errors"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/gsk148/gophkeeper/internal/app/server/storage"
	"github.com/gsk148/gophkeeper/internal/pkg/enc"
)

type UserService struct {
	db storage.IUserRepository
}

func NewUserService(db storage.IUserRepository) UserService {
	return UserService{db: db}
}

var ErrUserExists = errors.New("the user with specified name already exists")

func (s UserService) AddUser(ctx context.Context, req AuthReq) error {
	user := getUserFromRequest(req)
	userExist, err := doesUserExist(ctx, s.db, user)
	if err != nil {
		return err
	}
	if userExist {
		return ErrUserExists
	}

	hash, err := enc.HashPassword(user.Password)
	if err != nil {
		return nil
	}

	user.Password = hash
	_, err = s.db.AddUser(ctx, user)

	return err
}

func (s UserService) GetUser(ctx context.Context, user AuthReq) (storage.User, error) {
	su, err := s.db.GetUserByName(ctx, strings.ToLower(user.Name))
	if err != nil {
		return storage.User{}, err
	}

	if !enc.VerifyPassword(user.Password, su.Password) {
		log.Error(user.Password, su.Password)
		return storage.User{}, storage.ErrNotFound
	}

	return su, nil
}

func getUserFromRequest(r AuthReq) storage.User {
	return storage.User{
		Name:     strings.ToLower(r.Name),
		Password: r.Password,
	}
}

func doesUserExist(ctx context.Context, db storage.IUserRepository, user storage.User) (bool, error) {
	su, err := db.GetUserByName(ctx, user.Name)
	if err != nil && !errors.Is(err, storage.ErrNotFound) {
		return false, err
	}

	return su.ID != "", nil
}
