package handlers

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	log "github.com/sirupsen/logrus"

	"github.com/gsk148/gophkeeper/internal/app/server/models"
	"github.com/gsk148/gophkeeper/internal/app/server/services"
	"github.com/gsk148/gophkeeper/internal/pkg/services/data"
)

type IAuthService interface {
	Authorize(token string) (string, error)
	Login(ctx context.Context, cid string, user models.UserRequest) (string, string, error)
	Logout(ctx context.Context, cid string) (bool, error)
	Register(ctx context.Context, user models.UserRequest) error
}

type IBinaryService interface {
	DeleteBinary(ctx context.Context, uid, id string) error
	GetAllBinaries(ctx context.Context, uid string) ([]models.BinaryResponse, error)
	GetBinaryByID(ctx context.Context, uid, id string) (models.BinaryResponse, error)
	StoreBinary(ctx context.Context, uid string, data models.BinaryRequest) (string, error)
}

type ICardService interface {
	DeleteCard(ctx context.Context, uid, id string) error
	GetAllCards(ctx context.Context, uid string) ([]models.CardResponse, error)
	GetCardByID(ctx context.Context, uid, id string) (models.CardResponse, error)
	StoreCard(ctx context.Context, uid string, data models.CardRequest) (string, error)
}

type IPasswordService interface {
	DeletePassword(ctx context.Context, uid, id string) error
	GetAllPasswords(ctx context.Context, uid string) ([]models.PasswordResponse, error)
	GetPasswordByID(ctx context.Context, uid, id string) (models.PasswordResponse, error)
	StorePassword(ctx context.Context, uid string, data models.PasswordRequest) (string, error)
}

type ITextService interface {
	DeleteText(ctx context.Context, uid, id string) error
	GetAllTexts(ctx context.Context, uid string) ([]models.TextResponse, error)
	GetTextByID(ctx context.Context, uid, id string) (models.TextResponse, error)
	StoreText(ctx context.Context, uid string, data models.TextRequest) (string, error)
}

type Handler struct {
	authService     IAuthService
	binaryService   IBinaryService
	cardService     ICardService
	passwordService IPasswordService
	textService     ITextService
}

func NewHandler(repoURL string) (*chi.Mux, error) {
	h, err := initHandler(repoURL)
	if err != nil {
		return nil, err
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger, middleware.Compress(5, "/*"))

	r.Route("/api/v1/", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/login", h.Login())
			r.Post("/logout", h.Logout())
			r.Post("/register", h.Register())
		})

		r.With(h.Auth).Route("/storage", func(r chi.Router) {
			r.Route("/binary", func(r chi.Router) {
				r.Get("/", h.GetAllBinaries())
				r.Get("/{id}", h.GetBinaryByID())
				r.Post("/", h.StoreBinary())
				r.Delete("/{id}", h.DeleteBinary())
			})

			r.Route("/card", func(r chi.Router) {
				r.Get("/", h.GetAllCards())
				r.Get("/{id}", h.GetCardByID())
				r.Post("/", h.StoreCard())
				r.Delete("/{id}", h.DeleteCard())
			})

			r.Route("/password", func(r chi.Router) {
				r.Get("/", h.GetAllPasswords())
				r.Get("/{id}", h.GetPasswordByID())
				r.Post("/", h.StorePassword())
				r.Delete("/{id}", h.DeletePassword())
			})

			r.Route("/text", func(r chi.Router) {
				r.Get("/", h.GetAllTexts())
				r.Get("/{id}", h.GetTextByID())
				r.Post("/", h.StoreText())
				r.Delete("/{id}", h.DeleteText())
			})
		})
	})

	return r, nil
}

func initHandler(repoURL string) (Handler, error) {
	dataMS, err := data.NewService(repoURL)
	if err != nil {
		return Handler{}, err
	}

	authService, err := services.NewAuthService(repoURL)
	if err != nil {
		return Handler{}, err
	}

	return Handler{
		authService:     authService,
		binaryService:   services.NewBinaryService(dataMS),
		cardService:     services.NewCardService(dataMS),
		passwordService: services.NewPasswordService(dataMS),
		textService:     services.NewTextService(dataMS),
	}, nil
}

func handleHTTPError(w http.ResponseWriter, err error, code int) {
	log.Error(err)
	http.Error(w, http.StatusText(code), code)
}
