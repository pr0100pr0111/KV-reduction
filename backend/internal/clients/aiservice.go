package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/pr0100pr0111/KV-redaction/internal/models"
)

type AIServiceClient struct {
	BaseURL    string
	HTTPClient *http.Client
}

func NewAIServiceClient(baseURL string) *AIServiceClient {
	return &AIServiceClient{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 5 * time.Minute,
		},
	}
}

func (c *AIServiceClient) ProcessAudio(req models.AIServiceRequest) (*models.AIServiceResponse, error) {
	url := fmt.Sprintf("%s/process", c.BaseURL)

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result models.AIServiceResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *AIServiceClient) Health() error {
	resp, err := c.HTTPClient.Get(fmt.Sprintf("%s/health", c.BaseURL))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("AI service unhealthy: %s", resp.Status)
	}
	return nil
}
