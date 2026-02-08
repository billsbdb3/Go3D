package handlers

import (
	"3d-library/internal/models"
	"encoding/json"
	"net/http"

	"github.com/jmoiron/sqlx"
)

type SearchHandler struct {
	db *sqlx.DB
}

func NewSearchHandler(db *sqlx.DB) *SearchHandler {
	return &SearchHandler{db: db}
}

func (h *SearchHandler) Search(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "query parameter required", 400)
		return
	}

	var modelsList []models.Model
	err := h.db.Select(&modelsList, `
		SELECT DISTINCT m.* FROM models m
		LEFT JOIN model_files mf ON m.id = mf.model_id
		WHERE 
			m.name ILIKE $1 OR 
			m.description ILIKE $1 OR
			m.path ILIKE $1 OR
			mf.filename ILIKE $1
		ORDER BY m.created_at DESC
		LIMIT 100
	`, "%"+query+"%")

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	json.NewEncoder(w).Encode(modelsList)
}
