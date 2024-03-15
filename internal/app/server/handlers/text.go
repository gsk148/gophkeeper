package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/gsk148/gophkeeper/internal/app/models"
)

func (h Handler) DeleteText() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := r.Context().Value(uidKey).(string)
		id := chi.URLParam(r, "id")

		if err := h.textService.DeleteText(r.Context(), uid, id); err != nil {
			handleHTTPError(w, err, h.getErrorCode(err))
			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(""))
	}
}

func (h Handler) GetAllTexts() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := r.Context().Value(uidKey).(string)
		ts, err := h.textService.GetAllTexts(r.Context(), uid)
		if err != nil {
			handleHTTPError(w, err, h.getErrorCode(err))
			return
		}

		if err = json.NewEncoder(w).Encode(ts); err != nil {
			handleHTTPError(w, err, h.getErrorCode(err))
		}
	}
}

func (h Handler) GetTextByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := r.Context().Value(uidKey).(string)
		id := chi.URLParam(r, "id")

		t, err := h.textService.GetTextByID(r.Context(), uid, id)
		if err != nil {
			handleHTTPError(w, err, h.getErrorCode(err))
			return
		}

		if t.ID == "" {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusOK)
		}

		if err = json.NewEncoder(w).Encode(t); err != nil {
			handleHTTPError(w, err, h.getErrorCode(err))
		}
	}
}

func (h Handler) StoreText() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := r.Context().Value(uidKey).(string)

		var req models.TextRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			handleHTTPError(w, err, http.StatusBadRequest)
			return
		}

		id, err := h.textService.StoreText(r.Context(), uid, req)
		if err != nil {
			handleHTTPError(w, err, h.getErrorCode(err))
			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(id))
	}
}
