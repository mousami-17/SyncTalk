package cloudinary

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
	"time"

	"realtime-chat/src/config"
)

type CloudinaryResponse struct {
	PublicID     string `json:"public_id"`
	Version      int    `json:"version"`
	Signature    string `json:"signature"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
	Format       string `json:"format"`
	ResourceType string `json:"resource_type"`
	CreatedAt    string `json:"created_at"`
	Bytes        int64  `json:"bytes"`
	Type         string `json:"type"`
	URL          string `json:"url"`
	SecureURL    string `json:"secure_url"`
	Thumbnail    string `json:"thumbnail_url"`
	Duration     float64 `json:"duration"`
}

// UploadFile uploads a file to Cloudinary
func UploadFile(file multipart.File, filename string, resourceType string) (*CloudinaryResponse, error) {
	cloudName := config.Config.CloudinaryCloudName
	apiKey := config.Config.CloudinaryAPIKey
	apiSecret := config.Config.CloudinaryAPISecret

	if cloudName == "" || apiKey == "" || apiSecret == "" {
		return nil, fmt.Errorf("Cloudinary credentials not configured")
	}

	// Read file content
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}

	// Reset file pointer
	file.Seek(0, 0)

	// Prepare upload URL
	uploadURL := fmt.Sprintf("https://api.cloudinary.com/v1_1/%s/%s/upload", cloudName, resourceType)

	// Create multipart form
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add file
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return nil, err
	}
	part.Write(fileBytes)

	// Add parameters
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	writer.WriteField("timestamp", timestamp)
	writer.WriteField("api_key", apiKey)

	// Generate signature
	signature := generateSignature(timestamp, apiSecret)
	writer.WriteField("signature", signature)

	// Add optional parameters for better optimization
	if resourceType == "image" {
		writer.WriteField("quality", "auto")
		writer.WriteField("fetch_format", "auto")
	} else if resourceType == "video" {
		writer.WriteField("quality", "auto")
		writer.WriteField("resource_type", "video")
	}

	writer.Close()

	// Make request
	req, err := http.NewRequest("POST", uploadURL, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("upload request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("cloudinary error (status %d): %s", resp.StatusCode, string(bodyBytes))
	}

	// Parse response
	var cloudinaryResp CloudinaryResponse
	if err := json.NewDecoder(resp.Body).Decode(&cloudinaryResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	log.Printf("File uploaded to Cloudinary: %s", cloudinaryResp.SecureURL)
	return &cloudinaryResp, nil
}

// generateSignature creates a signature for Cloudinary API
func generateSignature(timestamp string, apiSecret string) string {
	// Create params string
	params := fmt.Sprintf("timestamp=%s%s", timestamp, apiSecret)
	
	// Generate SHA1 hash
	h := sha1.New()
	h.Write([]byte(params))
	return hex.EncodeToString(h.Sum(nil))
}

// GetResourceType determines the resource type from MIME type
func GetResourceType(mimeType string) string {
	if strings.HasPrefix(mimeType, "image/") {
		return "image"
	} else if strings.HasPrefix(mimeType, "video/") {
		return "video"
	} else if strings.HasPrefix(mimeType, "audio/") {
		return "video" // Cloudinary uses "video" for audio files
	}
	return "raw"
}

// GetFileType returns simplified file type
func GetFileType(mimeType string) string {
	if strings.HasPrefix(mimeType, "image/") {
		return "image"
	} else if strings.HasPrefix(mimeType, "video/") {
		return "video"
	} else if strings.HasPrefix(mimeType, "audio/") {
		return "audio"
	}
	return "file"
}

// GenerateThumbnailURL generates a thumbnail URL for images and videos
func GenerateThumbnailURL(publicID string, resourceType string) string {
	cloudName := config.Config.CloudinaryCloudName
	if resourceType == "image" {
		return fmt.Sprintf("https://res.cloudinary.com/%s/image/upload/w_300,h_300,c_fill/%s", cloudName, publicID)
	} else if resourceType == "video" {
		return fmt.Sprintf("https://res.cloudinary.com/%s/video/upload/w_300,h_300,c_fill,so_0/%s.jpg", cloudName, publicID)
	}
	return ""
}
