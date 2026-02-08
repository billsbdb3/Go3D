package handlers

import (
	"3d-library/internal/models"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

type TagHandler struct {
	db *sqlx.DB
}

func NewTagHandler(db *sqlx.DB) *TagHandler {
	return &TagHandler{db: db}
}

func (h *TagHandler) List(w http.ResponseWriter, r *http.Request) {
	var tags []models.Tag
	err := h.db.Select(&tags, "SELECT * FROM tags ORDER BY name")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	json.NewEncoder(w).Encode(tags)
}

func (h *TagHandler) AddToModel(w http.ResponseWriter, r *http.Request) {
	modelID := chi.URLParam(r, "id")
	var req struct {
		Tag string `json:"tag"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	// Get or create tag
	var tagID int64
	err := h.db.QueryRow(
		"INSERT INTO tags (name) VALUES ($1) ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name RETURNING id",
		req.Tag,
	).Scan(&tagID)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// Link to model
	_, err = h.db.Exec(
		"INSERT INTO model_tags (model_id, tag_id) VALUES ($1, $2) ON CONFLICT DO NOTHING",
		modelID, tagID,
	)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.WriteHeader(204)
}

func (h *TagHandler) GetModelTags(w http.ResponseWriter, r *http.Request) {
	modelID := chi.URLParam(r, "id")
	var tags []models.Tag
	err := h.db.Select(&tags, `
		SELECT t.* FROM tags t
		JOIN model_tags mt ON t.id = mt.tag_id
		WHERE mt.model_id = $1
		ORDER BY t.name
	`, modelID)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	json.NewEncoder(w).Encode(tags)
}
