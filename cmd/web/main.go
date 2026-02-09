package main

import (
	"3d-library/internal/database"
	"3d-library/internal/handlers"
	"3d-library/internal/jobs"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	db, err := database.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	log.Println("✓ Connected to database")

	jobClient := jobs.NewClient()
	defer jobClient.Close()

	log.Println("✓ Connected to Redis")

	// Initialize handlers
	libraryHandler := handlers.NewLibraryHandler(db)
	modelHandler := handlers.NewModelHandler(db)
	collectionHandler := handlers.NewCollectionHandler(db)
	tagHandler := handlers.NewTagHandler(db)
	fileHandler := handlers.NewFileHandler(db)
	scanHandler := handlers.NewScanHandler(db, jobClient)
	searchHandler := handlers.NewSearchHandler(db)
	uploadHandler := handlers.NewUploadHandler(db)

	// Setup router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	
	// CORS for development
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
			next.ServeHTTP(w, r)
		})
	})

	// Serve static files
	fileServer := http.FileServer(http.Dir("./web/static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fileServer))
	
	// Serve index.html at root
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./web/static/index.html")
	})

	// API routes
	r.Route("/api", func(r chi.Router) {
		// Libraries
		r.Get("/libraries", libraryHandler.List)
		r.Post("/libraries", libraryHandler.Create)
		r.Get("/libraries/{id}", libraryHandler.Get)
		r.Delete("/libraries/{id}", libraryHandler.Delete)
		r.Post("/libraries/{id}/scan", scanHandler.ScanLibrary)
		r.Post("/libraries/{id}/upload", uploadHandler.Upload)

		// Models
		r.Get("/models", modelHandler.List)
		r.Post("/models", modelHandler.Create)
		r.Get("/models/{id}", modelHandler.Get)
		r.Delete("/models/{id}", modelHandler.Delete)
		r.Get("/models/{id}/files", fileHandler.GetModelFiles)
		r.Post("/models/{id}/preview", modelHandler.SetPreview)
		r.Post("/models/{id}/tags", tagHandler.AddToModel)
		r.Get("/models/{id}/tags", tagHandler.GetModelTags)

		// Collections
		r.Get("/collections", collectionHandler.List)
		r.Post("/collections", collectionHandler.Create)
		r.Get("/collections/{id}", collectionHandler.Get)
		r.Get("/collections/{id}/models", collectionHandler.GetModels)
		r.Post("/collections/{id}/models", collectionHandler.AddModel)

		// Files
		r.Get("/files/{id}", fileHandler.Get)
		r.Get("/files/{id}/download", fileHandler.Serve)
		r.Delete("/files/{id}", fileHandler.Delete)

		// Tags
		r.Get("/tags", tagHandler.List)

		// Search
		r.Get("/search", searchHandler.Search)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("✓ Server running on http://192.168.3.26:%s", port)
	log.Println("✓ UI available at http://192.168.3.26:3000")
	log.Println("✓ Start worker with: go run cmd/worker/main.go")
	http.ListenAndServe(":"+port, r)
}
