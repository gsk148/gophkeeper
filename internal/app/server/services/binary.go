package services

import (
	"context"
	"errors"

	"github.com/gsk148/gophkeeper/internal/app/models"
	"github.com/gsk148/gophkeeper/internal/pkg/services/binary"
	"github.com/gsk148/gophkeeper/internal/pkg/services/data"
)

type BinaryService struct {
	binaryMS binary.Service
}

var (
	ErrBadArguments   = errors.New("the required arguments are not present")
	ErrBinaryNotFound = errors.New("requested binary data not found")
)

// NewBinaryService returns an instance of the BinaryService with pre-defined binary microservice.
func NewBinaryService(dataMS data.Service) *BinaryService {
	return &BinaryService{binaryMS: binary.NewService(dataMS)}
}

// DeleteBinary removes the stored data with the unique ID.
// The method removes the data of the specified user only.
func (s *BinaryService) DeleteBinary(ctx context.Context, uid, id string) error {
	if uid == "" || id == "" {
		return ErrBadArguments
	}
	err := s.binaryMS.DeleteBinary(ctx, uid, id)
	if errors.Is(err, binary.ErrNotFound) {
		return ErrBinaryNotFound
	}
	return err
}

// GetAllBinaries returns all the user's stored binaries.
func (s *BinaryService) GetAllBinaries(ctx context.Context, uid string) ([]models.BinaryResponse, error) {
	if uid == "" {
		return nil, ErrBadArguments
	}
	resp, err := s.binaryMS.GetAllBinaries(ctx, uid)
	if err != nil {
		if errors.Is(err, binary.ErrNotFound) {
			return nil, ErrBinaryNotFound
		}
		return nil, err
	}

	binaries := make([]models.BinaryResponse, 0, len(resp))
	for _, c := range resp {
		binaries = append(binaries, s.getResponseFromModel(c))
	}
	return binaries, nil
}

// GetBinaryByID returns the stored data by the unique ID.
// The method returns the data of the specified user only.
func (s *BinaryService) GetBinaryByID(ctx context.Context, uid, id string) (models.BinaryResponse, error) {
	if uid == "" || id == "" {
		return models.BinaryResponse{}, ErrBadArguments
	}
	resp, err := s.binaryMS.GetBinaryByID(ctx, uid, id)
	if err != nil {
		if errors.Is(err, binary.ErrNotFound) {
			return models.BinaryResponse{}, ErrBinaryNotFound
		}
		return models.BinaryResponse{}, err
	}
	return s.getResponseFromModel(resp), nil
}

// StoreBinary stores the original binary via the associated data microservice.
func (s *BinaryService) StoreBinary(ctx context.Context, uid string, binary models.BinaryRequest) (string, error) {
	if uid == "" || binary.Name == "" || binary.Data == nil {
		return "", ErrBadArguments
	}
	return s.binaryMS.StoreBinary(ctx, uid, s.getModelFromRequest(uid, binary))
}

func (s *BinaryService) getResponseFromModel(model binary.Binary) models.BinaryResponse {
	return models.BinaryResponse{
		UID:  model.UID,
		ID:   model.ID,
		Name: model.Name,
		Data: model.Data,
		Note: model.Note,
	}
}

func (s *BinaryService) getModelFromRequest(uid string, req models.BinaryRequest) binary.Binary {
	return binary.Binary{
		UID:  uid,
		Name: req.Name,
		Data: req.Data,
		Note: req.Note,
	}
}
