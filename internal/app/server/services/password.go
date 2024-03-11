package services

import (
	"context"
	"errors"

	"github.com/gsk148/gophkeeper/internal/app/models"
	"github.com/gsk148/gophkeeper/internal/pkg/services/data"
	"github.com/gsk148/gophkeeper/internal/pkg/services/password"
)

var ErrPasswordNotFound = errors.New("requested password data not found")

type PasswordService struct {
	passwordMS password.Service
}

// NewPasswordService returns an instance of the BinaryService with pre-defined password microservice.
func NewPasswordService(dataMS data.Service) *PasswordService {
	return &PasswordService{passwordMS: password.NewService(dataMS)}
}

// DeletePassword removes the stored data with the unique ID.
// The method removes the data of the specified user only.
func (s *PasswordService) DeletePassword(ctx context.Context, uid, id string) error {
	if uid == "" || id == "" {
		return ErrBadArguments
	}
	err := s.passwordMS.DeletePassword(ctx, uid, id)
	if errors.Is(err, password.ErrNotFound) {
		return ErrPasswordNotFound
	}
	return err
}

// GetAllPasswords returns all the user's stored passwords.
func (s *PasswordService) GetAllPasswords(ctx context.Context, uid string) ([]models.PasswordResponse, error) {
	if uid == "" {
		return nil, ErrBadArguments
	}
	resp, err := s.passwordMS.GetAllPasswords(ctx, uid)
	if err != nil {
		if errors.Is(err, password.ErrNotFound) {
			return nil, ErrPasswordNotFound
		}
		return nil, err
	}

	passwords := make([]models.PasswordResponse, 0, len(resp))
	for _, c := range resp {
		passwords = append(passwords, s.getResponseFromModel(c))
	}
	return passwords, nil
}

// GetPasswordByID returns the stored data by the unique ID.
// The method returns the data of the specified user only.
func (s *PasswordService) GetPasswordByID(ctx context.Context, uid, id string) (models.PasswordResponse, error) {
	if uid == "" || id == "" {
		return models.PasswordResponse{}, ErrBadArguments
	}
	res, err := s.passwordMS.GetPasswordByID(ctx, uid, id)
	if err != nil {
		if errors.Is(err, password.ErrNotFound) {
			return models.PasswordResponse{}, ErrPasswordNotFound
		}
		return models.PasswordResponse{}, err
	}
	return s.getResponseFromModel(res), nil
}

// StorePassword stores the original password via the associated data microservice.
func (s *PasswordService) StorePassword(ctx context.Context, uid string, req models.PasswordRequest) (string, error) {
	return s.passwordMS.StorePassword(ctx, s.getModelFromRequest(uid, req))
}

func (s *PasswordService) getResponseFromModel(model password.Password) models.PasswordResponse {
	return models.PasswordResponse{
		UID:      model.UID,
		ID:       model.ID,
		Name:     model.Name,
		User:     model.User,
		Password: model.Password,
		Note:     model.Note,
	}
}

func (s *PasswordService) getModelFromRequest(uid string, req models.PasswordRequest) password.Password {
	return password.Password{
		UID:      uid,
		Name:     req.Name,
		User:     req.User,
		Password: req.Password,
		Note:     req.Note,
	}
}
