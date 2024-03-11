package services

import (
	"context"
	"errors"

	"github.com/gsk148/gophkeeper/internal/app/models"
	"github.com/gsk148/gophkeeper/internal/pkg/services/card"
	"github.com/gsk148/gophkeeper/internal/pkg/services/data"
)

type CardService struct {
	cardMS card.Service
}

var ErrCardNotFound = errors.New("requested card data not found")

// NewCardService returns an instance of the BinaryService with pre-defined card microservice.
func NewCardService(dataMS data.Service) *CardService {
	return &CardService{cardMS: card.NewService(dataMS)}
}

// DeleteCard removes the stored data with the unique ID.
// The method removes the data of the specified user only.
func (s *CardService) DeleteCard(ctx context.Context, uid, id string) error {
	if uid == "" || id == "" {
		return ErrBadArguments
	}
	err := s.cardMS.DeleteCard(ctx, uid, id)
	if errors.Is(err, card.ErrNotFound) {
		return ErrCardNotFound
	}
	return err
}

// GetAllCards returns all the user's stored cards.
func (s *CardService) GetAllCards(ctx context.Context, uid string) ([]models.CardResponse, error) {
	if uid == "" {
		return nil, ErrBadArguments
	}
	resp, err := s.cardMS.GetAllCards(ctx, uid)
	if err != nil {
		if errors.Is(err, card.ErrNotFound) {
			return nil, ErrCardNotFound
		}
		return nil, err
	}

	cards := make([]models.CardResponse, 0, len(resp))
	for _, c := range resp {
		cards = append(cards, s.getResponseFromModel(c))
	}
	return cards, nil
}

// GetCardByID returns the stored data by the unique ID.
// The method returns the data of the specified user only.
func (s *CardService) GetCardByID(ctx context.Context, uid, id string) (models.CardResponse, error) {
	if uid == "" || id == "" {
		return models.CardResponse{}, ErrBadArguments
	}
	res, err := s.cardMS.GetCardByID(ctx, uid, id)
	if err != nil {
		if errors.Is(err, card.ErrNotFound) {
			return models.CardResponse{}, ErrCardNotFound
		}
		return models.CardResponse{}, err
	}
	return s.getResponseFromModel(res), nil
}

// StoreCard stores the original card via the associated data microservice.
func (s *CardService) StoreCard(ctx context.Context, uid string, card models.CardRequest) (string, error) {
	return s.cardMS.StoreCard(ctx, s.getModelFromRequest(uid, card))
}

func (s *CardService) getResponseFromModel(model card.Card) models.CardResponse {
	return models.CardResponse{
		UID:     model.UID,
		ID:      model.ID,
		Name:    model.Name,
		Number:  model.Number,
		Holder:  model.Holder,
		ExpDate: model.ExpDate,
		CVV:     model.CVV,
		Note:    model.Note,
	}
}

func (s *CardService) getModelFromRequest(uid string, req models.CardRequest) card.Card {
	return card.Card{
		UID:     uid,
		Name:    req.Name,
		Number:  req.Number,
		Holder:  req.Holder,
		ExpDate: req.ExpDate,
		CVV:     req.CVV,
		Note:    req.Note,
	}
}
