package storage

import (
	"context"
	"errors"
	"sync"

	"github.com/google/uuid"
)

type BasicStorage struct {
	data   sync.Map
	tokens sync.Map
	users  sync.Map
}

type DataStorage struct {
	user *sync.Map
}

func NewBasicStorage() *BasicStorage {
	return &BasicStorage{
		data:   sync.Map{},
		tokens: sync.Map{},
		users:  sync.Map{},
	}
}

func (s *BasicStorage) DeleteSession(ctx context.Context, cid string) error {
	s.tokens.Delete(cid)
	return nil
}

func (s *BasicStorage) GetSession(ctx context.Context, cid string) (string, error) {
	if t, ok := s.tokens.Load(cid); ok {
		return t.(string), nil
	}

	return "", errors.New("token not found")
}

func (s *BasicStorage) StoreSession(ctx context.Context, cid, token string) error {
	s.tokens.Store(cid, token)

	return nil
}

func (s *BasicStorage) AddUser(ctx context.Context, u User) (User, error) {
	id := uuid.NewString()
	u.ID = id
	s.users.Store(id, u)

	return u, nil
}

func (s *BasicStorage) GetUserByID(ctx context.Context, uid string) (User, error) {
	if u, ok := s.users.Load(uid); ok {
		return u.(User), nil
	}

	return User{}, ErrNotFound
}

func (s *BasicStorage) GetUserByName(ctx context.Context, name string) (User, error) {
	var user User

	s.users.Range(func(_, v any) bool {
		u := v.(User)
		if u.Name == name {
			user = u
			return false
		}

		return true
	})

	if user.ID == "" {
		return User{}, ErrNotFound
	}

	return user, nil
}

func (s *BasicStorage) DeleteUser(ctx context.Context, uid string) error {
	s.users.Delete(uid)

	return nil
}

func (s *BasicStorage) DeleteData(ctx context.Context, uid, id string) error {
	if d, ok := s.data.Load(id); ok {
		data := d.(SecureData)
		if data.UID == uid {
			s.data.Delete(id)
		}
	}

	return nil
}

func (s *BasicStorage) GetAllDataByType(ctx context.Context, uid string, t Type) ([]SecureData, error) {
	var data []SecureData
	if us, ok := s.data.Load(uid); ok {
		us.(DataStorage).user.Range(func(_, v any) bool {
			d := v.(SecureData)
			if d.Type == t {
				data = append(data, d)
			}
			return true
		})
	}

	return data, nil
}

func (s *BasicStorage) GetDataByID(ctx context.Context, uid, id string) (SecureData, error) {
	var (
		us any
		d  any
		ok bool
	)

	if us, ok = s.data.Load(uid); ok {
		if d, ok = us.(DataStorage).user.Load(id); ok {
			data := d.(SecureData)
			return data, nil
		}
	}
	return SecureData{}, errors.New("data not found")
}

func (s *BasicStorage) StoreData(ctx context.Context, data SecureData) (string, error) {
	id := uuid.NewString()
	data.ID = id

	if us, ok := s.data.Load(data.UID); !ok {
		sd := &sync.Map{}
		sd.Store(id, data)
		s.data.Store(data.UID, DataStorage{user: sd})
	} else {
		us.(DataStorage).user.Store(id, data)
	}

	return id, nil
}
