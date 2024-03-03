package services

import (
	"context"
	"github.com/segmentio/ksuid"

	"github.com/gsk148/gophkeeper/internal/app/server/storage"
	"github.com/gsk148/gophkeeper/internal/pkg/jwt"
)

type SessionService struct {
	db storage.ISessionRepository
}

func NewSessionService(db storage.ISessionRepository) SessionService {
	return SessionService{db: db}
}

func (s SessionService) RestoreSession(ctx context.Context, cid string) (string, error) {
	t, err := s.db.GetSession(ctx, cid)
	if err != nil {
		return "", err
	}

	if exp, eErr := jwt.IsTokenExpired(t); eErr != nil || exp {
		_ = s.DeleteSession(ctx, cid)
		if eErr != nil {
			return "", eErr
		}
		return "", jwt.ErrTokenExpired
	}

	return t, nil
}

func (s SessionService) StoreSession(ctx context.Context, token string) (string, error) {
	cid := generateClientID()
	return cid, s.db.StoreSession(ctx, cid, token)
}

func (s SessionService) DeleteSession(ctx context.Context, cid string) error {
	return s.db.DeleteSession(ctx, cid)
}

func (s SessionService) GenerateToken(uid string) (string, error) {
	return jwt.EncodeToken(uid)
}

func (s SessionService) GetUidFromToken(token string) (string, error) {
	return jwt.GetUserIDFromToken(token)
}

func (s SessionService) IsTokenExpired(token string) (bool, error) {
	if exp, err := jwt.IsTokenExpired(token); err != nil || exp {
		if err != nil {
			return true, err
		}
		return true, jwt.ErrTokenExpired
	}
	return false, nil
}

func generateClientID() string {
	return ksuid.New().String()
}
