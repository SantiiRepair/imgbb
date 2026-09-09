package uploader

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const ImgBBApiURL = "https://api.imgbb.com/1/upload"

type ImgBBResponse struct {
	Data struct {
		URL       string `json:"url"`
		DeleteURL string `json:"delete_url"`
	} `json:"data"`
	Success bool `json:"success"`
}

// Upload sends the image (in bytes) to ImgBB and returns the public URL
func Upload(ctx context.Context, apiKey string, imageBytes []byte) (string, error) {
	encoded := base64.StdEncoding.EncodeToString(imageBytes)

	data := url.Values{}
	data.Set("key", apiKey)
	data.Set("image", encoded)

	req, err := http.NewRequestWithContext(ctx, "POST", ImgBBApiURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{
		Timeout: 20 * time.Second,
	}
	
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result ImgBBResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if !result.Success {
		return "", fmt.Errorf("failed to upload to ImgBB")
	}

	return result.Data.URL, nil
}
