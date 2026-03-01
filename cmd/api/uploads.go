package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/google/uuid"
)

var allowedImageTypes = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
	"image/webp": ".webp",
}

func (app *application) uploadImageHandler(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 10<<20) // 10 MB

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		writeError(w, http.StatusBadRequest, "request too large or invalid multipart form")
		return
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		writeError(w, http.StatusBadRequest, "image field is required")
		return
	}
	defer file.Close()

	contentType := header.Header.Get("Content-Type")
	ext, ok := allowedImageTypes[contentType]
	if !ok {
		writeError(w, http.StatusUnsupportedMediaType, "only jpeg, png, and webp images are allowed")
		return
	}

	var imageURL string

	if app.config.cloudinary.cloudName != "" {
		// Production: upload to Cloudinary
		url, err := app.uploadToCloudinary(file, ext)
		if err != nil {
			app.logger.Printf("cloudinary upload: %v", err)
			writeError(w, http.StatusInternalServerError, "could not save image")
			return
		}
		imageURL = url
	} else {
		// Dev: save to local filesystem
		filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
		dst := filepath.Join(app.config.uploadsDir, filename)

		out, err := os.Create(dst)
		if err != nil {
			app.logger.Printf("create upload file: %v", err)
			writeError(w, http.StatusInternalServerError, "could not save image")
			return
		}
		defer out.Close()

		if _, err := io.Copy(out, file); err != nil {
			app.logger.Printf("write upload file: %v", err)
			writeError(w, http.StatusInternalServerError, "could not save image")
			return
		}
		imageURL = "/uploads/" + filename
	}

	writeJSON(w, http.StatusCreated, map[string]string{"url": imageURL}) //nolint:errcheck
}

// uploadToCloudinary signs and POSTs the image to Cloudinary's upload API using
// only the standard library — no external SDK required.
func (app *application) uploadToCloudinary(file io.Reader, ext string) (string, error) {
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)

	// Cloudinary signed upload: SHA-1("timestamp=<ts><api_secret>")
	h := sha1.New()
	h.Write([]byte("timestamp=" + timestamp + app.config.cloudinary.apiSecret))
	signature := hex.EncodeToString(h.Sum(nil))

	var body bytes.Buffer
	w := multipart.NewWriter(&body)
	_ = w.WriteField("api_key", app.config.cloudinary.apiKey)
	_ = w.WriteField("timestamp", timestamp)
	_ = w.WriteField("signature", signature)

	fw, err := w.CreateFormFile("file", "image"+ext)
	if err != nil {
		return "", fmt.Errorf("create form file: %w", err)
	}
	if _, err := io.Copy(fw, file); err != nil {
		return "", fmt.Errorf("copy file data: %w", err)
	}
	w.Close()

	endpoint := fmt.Sprintf(
		"https://api.cloudinary.com/v1_1/%s/image/upload",
		app.config.cloudinary.cloudName,
	)
	resp, err := http.Post(endpoint, w.FormDataContentType(), &body)
	if err != nil {
		return "", fmt.Errorf("http post: %w", err)
	}
	defer resp.Body.Close()

	var result struct {
		SecureURL string `json:"secure_url"`
		Error     struct {
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("decode response: %w", err)
	}
	if result.Error.Message != "" {
		return "", fmt.Errorf("cloudinary: %s", result.Error.Message)
	}
	return result.SecureURL, nil
}
