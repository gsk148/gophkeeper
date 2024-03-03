package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/gsk148/gophkeeper/internal/app/server/services"
	"github.com/gsk148/gophkeeper/internal/app/server/storage"
)

func (h Handler) GetAllPasswords() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := r.Context().Value("uid").(string)
		ps, err := services.GetAllPasswords(r.Context(), h.db, uid)
		if err != nil {
			handleHTTPError(w, err, http.StatusInternalServerError)
			return
		}

		if err = json.NewEncoder(w).Encode(ps); err != nil {
			handleHTTPError(w, err, http.StatusInternalServerError)
		}
	}
}

func (h Handler) GetPasswordByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := r.Context().Value("uid").(string)
		id := chi.URLParam(r, "id")

		p, err := services.GetPasswordByID(r.Context(), h.db, uid, id)
		if err != nil && errors.Is(err, storage.ErrNotFound) {
			handleHTTPError(w, err, http.StatusInternalServerError)
			return
		}

		if p.ID == "" {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusOK)
		}

		if err = json.NewEncoder(w).Encode(p); err != nil {
			handleHTTPError(w, err, http.StatusInternalServerError)
		}
	}
}

func (h Handler) StorePassword() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := r.Context().Value("uid").(string)

		var req services.PasswordReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			handleHTTPError(w, err, http.StatusBadRequest)
			return
		}

		id, err := services.StorePassword(r.Context(), h.db, uid, req)
		if err != nil {
			handleHTTPError(w, err, http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(id))
	}
}
