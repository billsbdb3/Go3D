package handlers

import (
	"3d-library/internal/models"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

type CollectionHandler struct {
	db *sqlx.DB
}

func NewCollectionHandler(db *sqlx.DB) *CollectionHandler {
	return &CollectionHandler{db: db}
}

func (h *CollectionHandler) List(w http.ResponseWriter, r *http.Request) {
	var collections []models.Collection
	err := h.db.Select(&collections, "SELECT * FROM collections ORDER BY created_at DESC")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	json.NewEncoder(w).Encode(collections)
}

func (h *CollectionHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var collection models.Collection
	err := h.db.Get(&collection, "SELECT * FROM collections WHERE id = $1", id)
	if err != nil {
		http.Error(w, "Not found", 404)
		return
	}
	json.NewEncoder(w).Encode(collection)
}

func (h *CollectionHandler) Create(w http.ResponseWriter, r *http.Request) {
	var collection models.Collection
	if err := json.NewDecoder(r.Body).Decode(&collection); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	err := h.db.QueryRow(
		"INSERT INTO collections (name) VALUES ($1) RETURNING id, created_at",
		collection.Name,
	).Scan(&collection.ID, &collection.CreatedAt)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.WriteHeader(201)
	json.NewEncoder(w).Encode(collection)
}

func (h *CollectionHandler) AddModel(w http.ResponseWriter, r *http.Request) {
	collectionID := chi.URLParam(r, "id")
	var req struct {
		ModelID int64 `json:"model_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	_, err := h.db.Exec(
		"INSERT INTO model_collections (model_id, collection_id) VALUES ($1, $2) ON CONFLICT DO NOTHING",
		req.ModelID, collectionID,
	)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.WriteHeader(204)
}

func (h *CollectionHandler) GetModels(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var modelsList []models.Model
	err := h.db.Select(&modelsList, `
		SELECT m.* FROM models m
		JOIN model_collections mc ON m.id = mc.model_id
		WHERE mc.collection_id = $1
		ORDER BY m.created_at DESC
	`, id)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	json.NewEncoder(w).Encode(modelsList)
}
