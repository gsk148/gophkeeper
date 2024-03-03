package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/gsk148/gophkeeper/internal/app/server/services"
	"github.com/gsk148/gophkeeper/internal/app/server/storage"
)

func (h Handler) GetAllBinaries() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := r.Context().Value("uid").(string)
		bs, err := services.GetAllBinaries(r.Context(), h.db, uid)
		if err != nil {
			handleHTTPError(w, err, http.StatusInternalServerError)
			return
		}

		if err = json.NewEncoder(w).Encode(bs); err != nil {
			handleHTTPError(w, err, http.StatusInternalServerError)
		}
	}
}

func (h Handler) GetBinaryByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := r.Context().Value("uid").(string)
		id := chi.URLParam(r, "id")

		b, err := services.GetBinaryByID(r.Context(), h.db, uid, id)
		if err != nil && errors.Is(err, storage.ErrNotFound) {
			handleHTTPError(w, err, http.StatusInternalServerError)
			return
		}

		if b.ID == "" {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusOK)
		}

		if err = json.NewEncoder(w).Encode(b); err != nil {
			handleHTTPError(w, err, http.StatusInternalServerError)
		}
	}
}

func (h Handler) StoreBinary() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := r.Context().Value("uid").(string)

		var req services.BinaryReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			handleHTTPError(w, err, http.StatusBadRequest)
			return
		}

		id, err := services.StoreBinary(r.Context(), h.db, uid, req)
		if err != nil {
			handleHTTPError(w, err, http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(id))
	}
}
