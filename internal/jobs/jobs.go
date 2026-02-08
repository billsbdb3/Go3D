package jobs

import (
	"3d-library/internal/scanner"
	"context"
	"encoding/json"
	"log"
	"path/filepath"

	"github.com/hibiken/asynq"
	"github.com/jmoiron/sqlx"
)

const (
	TypeScanLibrary = "library:scan"
)

type ScanLibraryPayload struct {
	LibraryID int64  `json:"library_id"`
	Path      string `json:"path"`
}

func NewScanLibraryTask(libraryID int64, path string) (*asynq.Task, error) {
	payload, err := json.Marshal(ScanLibraryPayload{LibraryID: libraryID, Path: path})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeScanLibrary, payload), nil
}

func HandleScanLibraryTask(ctx context.Context, t *asynq.Task, db *sqlx.DB) error {
	var p ScanLibraryPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return err
	}

	log.Printf("Scanning library %d at %s", p.LibraryID, p.Path)

	s := scanner.New(p.Path)
	files, err := s.Scan()
	if err != nil {
		return err
	}

	// Group files by directory
	modelDirs := make(map[string][]scanner.FileInfo)
	for _, file := range files {
		dir := filepath.Dir(file.Path)
		modelDirs[dir] = append(modelDirs[dir], file)
	}

	added := 0
	for modelPath, dirFiles := range modelDirs {
		modelName := filepath.Base(modelPath)
		
		var modelID int64
		err := db.QueryRow(`
			INSERT INTO models (library_id, name, path) 
			VALUES ($1, $2, $3) 
			ON CONFLICT (library_id, path) DO UPDATE 
			SET updated_at = NOW() 
			RETURNING id
		`, p.LibraryID, modelName, modelPath).Scan(&modelID)
		
		if err != nil {
			continue
		}

		for _, file := range dirFiles {
			_, err = db.Exec(`
				INSERT INTO model_files (model_id, filename, path, size, mime_type, digest) 
				VALUES ($1, $2, $3, $4, $5, $6) 
				ON CONFLICT (path) DO UPDATE SET size = $4, digest = $6
			`, modelID, filepath.Base(file.Path), file.Path, file.Size, file.MimeType, file.Digest)
			
			if err == nil {
				added++
			}
		}
	}

	log.Printf("Scan complete: %d files scanned, %d models, %d files added", len(files), len(modelDirs), added)
	return nil
}

func NewClient() *asynq.Client {
	return asynq.NewClient(asynq.RedisClientOpt{Addr: "localhost:6379"})
}

func NewServer(db *sqlx.DB) *asynq.ServeMux {
	mux := asynq.NewServeMux()
	mux.HandleFunc(TypeScanLibrary, func(ctx context.Context, t *asynq.Task) error {
		return HandleScanLibraryTask(ctx, t, db)
	})
	return mux
}
