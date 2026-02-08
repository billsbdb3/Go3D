package handlers

import (
	"3d-library/internal/models"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

type FileHandler struct {
	db *sqlx.DB
}

func NewFileHandler(db *sqlx.DB) *FileHandler {
	return &FileHandler{db: db}
}

func (h *FileHandler) GetModelFiles(w http.ResponseWriter, r *http.Request) {
	modelID := chi.URLParam(r, "id")
	var files []models.ModelFile
	err := h.db.Select(&files, "SELECT * FROM model_files WHERE model_id = $1 ORDER BY filename", modelID)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	json.NewEncoder(w).Encode(files)
}

func (h *FileHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var file models.ModelFile
	err := h.db.Get(&file, "SELECT * FROM model_files WHERE id = $1", id)
	if err != nil {
		http.Error(w, "Not found", 404)
		return
	}
	json.NewEncoder(w).Encode(file)
}

func (h *FileHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	_, err := h.db.Exec("DELETE FROM model_files WHERE id = $1", id)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.WriteHeader(204)
}

func (h *FileHandler) Serve(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var file models.ModelFile
	err := h.db.Get(&file, "SELECT * FROM model_files WHERE id = $1", id)
	if err != nil {
		http.Error(w, "Not found", 404)
		return
	}
	http.ServeFile(w, r, file.Path)
}
