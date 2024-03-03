package storage

import (
	"context"
	"errors"
)

var ErrNotFound = errors.New("data not found")

type Type int

const (
	SBinary Type = iota
	SCard
	SPassword
	SText
)

type IRepository interface {
	IDataRepository
	ISessionRepository
	IUserRepository
}

type IDataRepository interface {
	GetAllDataByType(ctx context.Context, uid string, t Type) ([]SecureData, error)
	GetDataByID(ctx context.Context, uid, id string) (SecureData, error)
	StoreData(ctx context.Context, data SecureData) (string, error)
	DeleteData(ctx context.Context, uid, id string) error
}

type ISessionRepository interface {
	DeleteSession(ctx context.Context, cid string) error
	GetSession(ctx context.Context, cid string) (string, error)
	StoreSession(ctx context.Context, cid, token string) error
}

type IUserRepository interface {
	AddUser(ctx context.Context, user User) (User, error)
	DeleteUser(ctx context.Context, uid string) error
	GetUserByID(ctx context.Context, uid string) (User, error)
	GetUserByName(ctx context.Context, name string) (User, error)
}

type SecureData struct {
	UID  string `json:"-"`
	ID   string `json:"id"`
	Data []byte `json:"data"`
	Type Type   `json:"-"`
}

type User struct {
	ID       string `json:"-"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

func NewStorage() (IRepository, error) {
	return NewBasicStorage(), nil
}
