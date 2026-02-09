package handlers

import (
	"3d-library/internal/models"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

type ModelHandler struct {
	db *sqlx.DB
}

func NewModelHandler(db *sqlx.DB) *ModelHandler {
	return &ModelHandler{db: db}
}

func (h *ModelHandler) List(w http.ResponseWriter, r *http.Request) {
	libraryID := r.URL.Query().Get("library_id")
	
	var modelsList []models.Model
	var err error
	
	if libraryID != "" {
		err = h.db.Select(&modelsList, "SELECT * FROM models WHERE library_id = $1 ORDER BY created_at DESC", libraryID)
	} else {
		err = h.db.Select(&modelsList, "SELECT * FROM models ORDER BY created_at DESC LIMIT 100")
	}
	
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	json.NewEncoder(w).Encode(modelsList)
}

func (h *ModelHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var model models.Model
	err := h.db.Get(&model, "SELECT * FROM models WHERE id = $1", id)
	if err != nil {
		http.Error(w, "Not found", 404)
		return
	}
	json.NewEncoder(w).Encode(model)
}

func (h *ModelHandler) Create(w http.ResponseWriter, r *http.Request) {
	var model models.Model
	if err := json.NewDecoder(r.Body).Decode(&model); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	// Use ON CONFLICT to handle duplicates
	err := h.db.QueryRow(`
		INSERT INTO models (library_id, name, path, description) 
		VALUES ($1, $2, $3, $4) 
		ON CONFLICT (library_id, path) DO UPDATE 
		SET name = EXCLUDED.name, description = EXCLUDED.description, updated_at = NOW()
		RETURNING id, created_at, updated_at
	`, model.LibraryID, model.Name, model.Path, model.Description).Scan(&model.ID, &model.CreatedAt, &model.UpdatedAt)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.WriteHeader(201)
	json.NewEncoder(w).Encode(model)
}

func (h *ModelHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	_, err := h.db.Exec("DELETE FROM models WHERE id = $1", id)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.WriteHeader(204)
}

func (h *ModelHandler) SetPreview(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req struct {
		FileID *int `json:"file_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	_, err := h.db.Exec("UPDATE models SET preview_file_id = $1 WHERE id = $2", req.FileID, id)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.WriteHeader(204)
}
