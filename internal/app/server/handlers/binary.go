package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/gsk148/gophkeeper/internal/app/models"
)

func (h Handler) DeleteBinary() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := r.Context().Value(uidKey).(string)
		id := chi.URLParam(r, "id")

		if err := h.binaryService.DeleteBinary(r.Context(), uid, id); err != nil {
			handleHTTPError(w, err, h.getErrorCode(err))
			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(""))
	}
}

func (h Handler) GetAllBinaries() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := r.Context().Value(uidKey).(string)
		bs, err := h.binaryService.GetAllBinaries(r.Context(), uid)
		if err != nil {
			handleHTTPError(w, err, h.getErrorCode(err))
			return
		}

		if err = json.NewEncoder(w).Encode(bs); err != nil {
			handleHTTPError(w, err, h.getErrorCode(err))
		}
	}
}

func (h Handler) GetBinaryByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := r.Context().Value(uidKey).(string)
		id := chi.URLParam(r, "id")

		b, err := h.binaryService.GetBinaryByID(r.Context(), uid, id)
		if err != nil {
			handleHTTPError(w, err, h.getErrorCode(err))
			return
		}

		if b.ID == "" {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusOK)
		}

		if err = json.NewEncoder(w).Encode(b); err != nil {
			handleHTTPError(w, err, h.getErrorCode(err))
		}
	}
}

func (h Handler) StoreBinary() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := r.Context().Value(uidKey).(string)

		var req models.BinaryRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			handleHTTPError(w, err, http.StatusBadRequest)
			return
		}

		id, err := h.binaryService.StoreBinary(r.Context(), uid, req)
		if err != nil {
			handleHTTPError(w, err, h.getErrorCode(err))
			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(id))
	}
}
