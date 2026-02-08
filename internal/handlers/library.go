package handlers

import (
	"3d-library/internal/models"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

type LibraryHandler struct {
	db *sqlx.DB
}

func NewLibraryHandler(db *sqlx.DB) *LibraryHandler {
	return &LibraryHandler{db: db}
}

func (h *LibraryHandler) List(w http.ResponseWriter, r *http.Request) {
	var libraries []models.Library
	err := h.db.Select(&libraries, "SELECT * FROM libraries ORDER BY created_at DESC")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	json.NewEncoder(w).Encode(libraries)
}

func (h *LibraryHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var library models.Library
	err := h.db.Get(&library, "SELECT * FROM libraries WHERE id = $1", id)
	if err != nil {
		http.Error(w, "Not found", 404)
		return
	}
	json.NewEncoder(w).Encode(library)
}

func (h *LibraryHandler) Create(w http.ResponseWriter, r *http.Request) {
	var library models.Library
	if err := json.NewDecoder(r.Body).Decode(&library); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	err := h.db.QueryRow(
		"INSERT INTO libraries (name, path, storage) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at",
		library.Name, library.Path, library.Storage,
	).Scan(&library.ID, &library.CreatedAt, &library.UpdatedAt)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.WriteHeader(201)
	json.NewEncoder(w).Encode(library)
}

func (h *LibraryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	_, err := h.db.Exec("DELETE FROM libraries WHERE id = $1", id)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.WriteHeader(204)
}
