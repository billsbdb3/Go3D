package handlers

import (
	"archive/zip"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

type UploadHandler struct {
	db *sqlx.DB
}

func NewUploadHandler(db *sqlx.DB) *UploadHandler {
	return &UploadHandler{db: db}
}

func (h *UploadHandler) Upload(w http.ResponseWriter, r *http.Request) {
	libraryID := chi.URLParam(r, "id")
	
	var library struct {
		ID   int64  `db:"id"`
		Path string `db:"path"`
	}
	err := h.db.Get(&library, "SELECT id, path FROM libraries WHERE id = $1", libraryID)
	if err != nil {
		http.Error(w, "Library not found", 404)
		return
	}

	err = r.ParseMultipartForm(100 << 20)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	modelName := r.FormValue("model_name")
	if modelName == "" {
		http.Error(w, "model_name required", 400)
		return
	}

	modelPath := filepath.Join(library.Path, modelName)
	os.MkdirAll(modelPath, 0755)

	files := r.MultipartForm.File["files[]"]
	if len(files) == 0 {
		files = r.MultipartForm.File["files"]
	}
	if len(files) == 0 {
		http.Error(w, "No files uploaded", 400)
		return
	}

	var modelID int64
	err = h.db.QueryRow(
		"INSERT INTO models (library_id, name, path) VALUES ($1, $2, $3) ON CONFLICT (library_id, path) DO UPDATE SET name = $2 RETURNING id",
		library.ID, modelName, modelPath,
	).Scan(&modelID)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	uploaded := []string{}
	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			continue
		}
		defer file.Close()

		if strings.HasSuffix(strings.ToLower(fileHeader.Filename), ".zip") {
			tmpZip := filepath.Join(os.TempDir(), fileHeader.Filename)
			tmpFile, err := os.Create(tmpZip)
			if err != nil {
				log.Printf("Error creating temp file: %v", err)
				continue
			}
			io.Copy(tmpFile, file)
			tmpFile.Close()

			extracted, err := h.extractZip(tmpZip, modelPath, modelID)
			if err != nil {
				log.Printf("Error extracting ZIP: %v", err)
			} else {
				uploaded = append(uploaded, extracted...)
			}
			os.Remove(tmpZip)
			continue
		}

		destPath := filepath.Join(modelPath, fileHeader.Filename)
		destFile, err := os.Create(destPath)
		if err != nil {
			continue
		}

		hash := sha256.New()
		size, _ := io.Copy(io.MultiWriter(destFile, hash), file)
		destFile.Close()

		digest := fmt.Sprintf("%x", hash.Sum(nil))

		_, err = h.db.Exec(
			"INSERT INTO model_files (model_id, filename, path, size, digest) VALUES ($1, $2, $3, $4, $5) ON CONFLICT (path) DO UPDATE SET size = EXCLUDED.size, digest = EXCLUDED.digest",
			modelID, fileHeader.Filename, destPath, size, digest,
		)
		if err != nil {
			log.Printf("Error saving file to DB: %v", err)
		} else {
			uploaded = append(uploaded, fileHeader.Filename)
		}
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"uploaded": uploaded,
		"count":    len(uploaded),
	})
}

func (h *UploadHandler) extractZip(zipPath, destDir string, modelID int64) ([]string, error) {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	extracted := []string{}
	for _, f := range r.File {
		if f.FileInfo().IsDir() || strings.HasPrefix(filepath.Base(f.Name), ".") {
			continue
		}

		fpath := filepath.Join(destDir, f.Name)
		os.MkdirAll(filepath.Dir(fpath), 0755)

		outFile, err := os.Create(fpath)
		if err != nil {
			log.Printf("Error creating file %s: %v", fpath, err)
			continue
		}

		rc, err := f.Open()
		if err != nil {
			outFile.Close()
			log.Printf("Error opening ZIP entry %s: %v", f.Name, err)
			continue
		}

		hash := sha256.New()
		size, _ := io.Copy(io.MultiWriter(outFile, hash), rc)
		outFile.Close()
		rc.Close()

		digest := fmt.Sprintf("%x", hash.Sum(nil))

		_, err = h.db.Exec(
			"INSERT INTO model_files (model_id, filename, path, size, digest) VALUES ($1, $2, $3, $4, $5) ON CONFLICT (path) DO UPDATE SET size = EXCLUDED.size, digest = EXCLUDED.digest",
			modelID, filepath.Base(f.Name), fpath, size, digest,
		)
		if err != nil {
			log.Printf("Error saving %s to DB: %v", f.Name, err)
		} else {
			extracted = append(extracted, f.Name)
		}
	}

	return extracted, nil
}
