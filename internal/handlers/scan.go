package handlers

import (
	"3d-library/internal/jobs"
	"3d-library/internal/models"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/hibiken/asynq"
	"github.com/jmoiron/sqlx"
)

type ScanHandler struct {
	db     *sqlx.DB
	client *asynq.Client
}

func NewScanHandler(db *sqlx.DB, client *asynq.Client) *ScanHandler {
	return &ScanHandler{db: db, client: client}
}

func (h *ScanHandler) ScanLibrary(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	
	var library models.Library
	err := h.db.Get(&library, "SELECT * FROM libraries WHERE id = $1", id)
	if err != nil {
		http.Error(w, "Library not found", 404)
		return
	}

	// Queue the scan job
	task, err := jobs.NewScanLibraryTask(library.ID, library.Path)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	info, err := h.client.Enqueue(task)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Scan queued",
		"job_id":  info.ID,
	})
}
