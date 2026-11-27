package imagekit

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"time"
)

var (
	uploadEndpoint = "https://upload.imagekit.io/api/v1/files/upload"
	client         = &http.Client{Timeout: 30 * time.Second}
)

type uploadResponse struct {
	URL     string `json:"url"`
	Message string `json:"message"`
}

func Upload(ctx context.Context, privateKey, folder string, content []byte, filename string) (string, error) {
	var body bytes.Buffer
	form := multipart.NewWriter(&body)
	for key, value := range map[string]string{"fileName": filename, "folder": folder, "useUniqueFileName": "true"} {
		if err := form.WriteField(key, value); err != nil {
			return "", err
		}
	}
	part, err := form.CreateFormFile("file", filename)
	if err != nil {
		return "", err
	}
	if _, err := part.Write(content); err != nil {
		return "", err
	}
	if err := form.Close(); err != nil {
		return "", err
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, uploadEndpoint, &body)
	if err != nil {
		return "", err
	}
	request.SetBasicAuth(privateKey, "")
	request.Header.Set("Content-Type", form.FormDataContentType())
	response, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	var result uploadResponse
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("imagekit: status %d", response.StatusCode)
	}
	if response.StatusCode != http.StatusOK || result.URL == "" {
		return "", fmt.Errorf("imagekit: status %d: %s", response.StatusCode, result.Message)
	}
	return result.URL, nil
}
