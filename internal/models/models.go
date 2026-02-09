package models

import "time"

type Library struct {
	ID        int64     `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	Path      string    `db:"path" json:"path"`
	Storage   string    `db:"storage" json:"storage"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type Model struct {
	ID            int64     `db:"id" json:"id"`
	LibraryID     int64     `db:"library_id" json:"library_id"`
	Name          string    `db:"name" json:"name"`
	Path          string    `db:"path" json:"path"`
	Description   *string   `db:"description" json:"description"`
	PreviewFileID *int64    `db:"preview_file_id" json:"preview_file_id"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at"`
}

type ModelFile struct {
	ID        int64     `db:"id" json:"id"`
	ModelID   int64     `db:"model_id" json:"model_id"`
	Filename  string    `db:"filename" json:"filename"`
	Path      string    `db:"path" json:"path"`
	Size      int64     `db:"size" json:"size"`
	MimeType  *string   `db:"mime_type" json:"mime_type"`
	Digest    *string   `db:"digest" json:"digest"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type Collection struct {
	ID        int64     `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type Tag struct {
	ID   int64  `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}
