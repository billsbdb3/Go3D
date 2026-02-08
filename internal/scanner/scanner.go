package scanner

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type FileInfo struct {
	Path     string
	Size     int64
	Digest   string
	MimeType string
}

type Scanner struct {
	rootPath string
}

func New(rootPath string) *Scanner {
	return &Scanner{rootPath: rootPath}
}

func (s *Scanner) Scan() ([]FileInfo, error) {
	var files []FileInfo

	err := filepath.Walk(s.rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		// Only scan 3D model files
		ext := strings.ToLower(filepath.Ext(path))
		if !is3DFile(ext) {
			return nil
		}

		digest, err := calculateDigest(path)
		if err != nil {
			return err
		}

		files = append(files, FileInfo{
			Path:     path,
			Size:     info.Size(),
			Digest:   digest,
			MimeType: getMimeType(ext),
		})

		return nil
	})

	return files, err
}

func is3DFile(ext string) bool {
	supported := []string{".stl", ".obj", ".3mf", ".ply", ".gcode"}
	for _, s := range supported {
		if ext == s {
			return true
		}
	}
	return false
}

func getMimeType(ext string) string {
	types := map[string]string{
		".stl":   "model/stl",
		".obj":   "model/obj",
		".3mf":   "model/3mf",
		".ply":   "model/ply",
		".gcode": "text/x-gcode",
	}
	return types[ext]
}

func calculateDigest(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
