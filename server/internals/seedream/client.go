package seedream

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

// defaultAPIKey is used only if no override is provided and SEEDREAM_API_KEY is not set.
// NOTE: Hardcoding secrets is not recommended for production use.
const defaultAPIKey = "2d2a1b66dc9955688223ae64e015a1b3"

type generateRequest struct {
	Prompt          string   `json:"prompt"`
	ImageSize       string   `json:"image_size,omitempty"`
	ImageResolution string   `json:"image_resolution,omitempty"`
	MaxImages       int      `json:"max_images,omitempty"`
	Seed            int      `json:"seed,omitempty"`
	ReferenceImages []string `json:"reference_images,omitempty"` // base64 data URLs
}

type GenerateResponse struct {
	JobID  string   `json:"job_id,omitempty"`
	Status string   `json:"status,omitempty"`
	Images []string `json:"images,omitempty"`
	Error  string   `json:"error,omitempty"`
}

// GenerateImages calls the Seedream API with 0/1/2 input images (as raw bytes).
// The API key is read from SEEDREAM_API_KEY env var unless an override is provided.
func GenerateImages(ctx context.Context, apiKeyOverride string, prompt string, imageBytes [][]byte) (*GenerateResponse, error) {
	apiKey := apiKeyOverride
	if apiKey == "" {
		apiKey = os.Getenv("SEEDREAM_API_KEY")
	}
	if apiKey == "" {
		apiKey = defaultAPIKey
	}

	// Pre-allocate slice for reference images
	refImgs := make([]string, 0, len(imageBytes))
	for _, b := range imageBytes {
		if len(b) == 0 {
			continue
		}
		// Optimize base64 encoding with buffer
		encoded := base64.StdEncoding.EncodeToString(b)
		refImgs = append(refImgs, "data:application/octet-stream;base64,"+encoded)
	}

	reqBody := generateRequest{
		Prompt:          prompt,
		ImageSize:       "portrait_3_4",
		ImageResolution: "4K",
		MaxImages:       1,
		Seed:            54321,
	}
	if len(refImgs) > 0 {
		reqBody.ReferenceImages = refImgs
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", "https://api.kie.ai/api/v1/jobs/createTask", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Authorization", "Bearer "+apiKey)
	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 90 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        50,
			MaxIdleConnsPerHost: 10,
			IdleConnTimeout:     90 * time.Second,
			ForceAttemptHTTP2:   true,
		},
	}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var apiResp GenerateResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		if apiResp.Error != "" {
			return nil, fmt.Errorf("seedream error: %s", apiResp.Error)
		}
		return nil, fmt.Errorf("seedream bad status: %d", resp.StatusCode)
	}
	return &apiResp, nil
}
