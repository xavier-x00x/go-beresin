package utils

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Allowed MIME types grouped by category.
var allowedMIMETypes = map[string]bool{
	// Photos
	"image/jpeg": true,
	"image/jpg":  true,
	"image/png":  true,
	"image/webp": true,
	"image/gif":  true,

	// Videos
	"video/mp4":  true,
	"video/mpeg": true,
	"video/quicktime": true, // .mov

	// Documents
	"application/pdf":  true,
	"application/msword": true,
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true, // .docx
}

const (
	// MaxFileSize specifies maximum upload size of 10MB
	MaxFileSize = 10 * 1024 * 1024
)

// ProcessUploadedFile validates and saves a multipart file upload to a target directory.
// Returns the relative path of the saved file or an error.
func ProcessUploadedFile(fileHeader *multipart.FileHeader, destDir string) (string, error) {
	// 1. Validate File Size
	if fileHeader.Size > MaxFileSize {
		return "", fmt.Errorf("file size (%d bytes) exceeds maximum limit of %d MB", fileHeader.Size, MaxFileSize/(1024*1024))
	}

	// Open the file
	file, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer file.Close()

	// 2. Validate MIME Type by sniffing the first 512 bytes
	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return "", fmt.Errorf("failed to read file content for verification: %w", err)
	}

	// Reset file read pointer after reading sniff buffer
	if _, err := file.Seek(0, 0); err != nil {
		return "", fmt.Errorf("failed to reset file pointer: %w", err)
	}

	detectedMIME := http.DetectContentType(buffer[:n])
	// Split MIME parameters if any (e.g. text/html; charset=utf-8)
	detectedMIME = strings.Split(detectedMIME, ";")[0]

	// Check against our allowed list
	if !allowedMIMETypes[detectedMIME] {
		return "", fmt.Errorf("unsupported file type: %s", detectedMIME)
	}

	// 3. Ensure the destination directory exists
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create upload destination directory: %w", err)
	}

	// 4. Generate a unique, safe filename to prevent collisions
	ext := filepath.Ext(fileHeader.Filename)
	if ext == "" {
		// Fallback extensions based on MIME if extension is missing
		switch detectedMIME {
		case "image/jpeg", "image/jpg":
			ext = ".jpg"
		case "image/png":
			ext = ".png"
		case "image/webp":
			ext = ".webp"
		case "video/mp4":
			ext = ".mp4"
		case "application/pdf":
			ext = ".pdf"
		}
	}

	// Clean the original filename (sanitize it) and prepend a UUID
	uniqueID := uuid.New().String()
	sanitizedName := sanitizeFilename(filepath.Base(fileHeader.Filename))
	// Truncate name to avoid too long paths
	if len(sanitizedName) > 50 {
		sanitizedName = sanitizedName[:50]
	}

	filename := fmt.Sprintf("%d_%s_%s%s", time.Now().Unix(), uniqueID[:8], sanitizedName, ext)
	filePath := filepath.Join(destDir, filename)

	// Create target file
	out, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create target file on disk: %w", err)
	}
	defer out.Close()

	// Copy content
	if _, err = io.Copy(out, file); err != nil {
		return "", fmt.Errorf("failed to save file content: %w", err)
	}

	return filePath, nil
}

// sanitizeFilename replaces non-alphanumeric characters with underscores.
func sanitizeFilename(name string) string {
	ext := filepath.Ext(name)
	base := strings.TrimSuffix(name, ext)

	var builder strings.Builder
	for _, r := range base {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '-' || r == '_' {
			builder.WriteRune(r)
		} else {
			builder.WriteRune('_')
		}
	}
	return builder.String()
}
