package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/gsk148/gophkeeper/internal/app/server/models"
	"github.com/gsk148/gophkeeper/internal/app/server/services"
)

func (h Handler) DeleteCard() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := r.Context().Value(uidKey).(string)
		id := chi.URLParam(r, "id")

		if err := h.cardService.DeleteCard(r.Context(), uid, id); err != nil {
			handleHTTPError(w, err, http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(""))
	}
}

func (h Handler) GetAllCards() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := r.Context().Value(uidKey).(string)
		cs, err := h.cardService.GetAllCards(r.Context(), uid)
		if err != nil {
			handleHTTPError(w, err, http.StatusInternalServerError)
			return
		}

		if err = json.NewEncoder(w).Encode(cs); err != nil {
			handleHTTPError(w, err, http.StatusInternalServerError)
		}
	}
}

func (h Handler) GetCardByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := r.Context().Value(uidKey).(string)
		id := chi.URLParam(r, "id")

		c, err := h.cardService.GetCardByID(r.Context(), uid, id)
		if err != nil && errors.Is(err, services.ErrCardNotFound) {
			handleHTTPError(w, err, http.StatusInternalServerError)
			return
		}

		if c.ID == "" {
			handleHTTPError(w, services.ErrCardNotFound, http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
		if err = json.NewEncoder(w).Encode(c); err != nil {
			handleHTTPError(w, err, http.StatusInternalServerError)
		}
	}
}

func (h Handler) StoreCard() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := r.Context().Value(uidKey).(string)

		var req models.CardRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			handleHTTPError(w, err, http.StatusBadRequest)
			return
		}

		id, err := h.cardService.StoreCard(r.Context(), uid, req)
		if err != nil {
			handleHTTPError(w, err, http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(id))
	}
}
