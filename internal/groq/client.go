package groq

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"groq-cli-chat/resources"
)

type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

func NewClient(baseURL, apiKey string) (*Client, error) {
	if baseURL == "" || apiKey == "" {
		return nil, fmt.Errorf(resources.ErrInvalidClientParams)
	}
	return &Client{
		baseURL:    baseURL,
		apiKey:     apiKey,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}, nil
}

func (c *Client) ListModels() ([]string, error) {
	resp, err := c.makeRequest("GET", "models", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Data []struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf(resources.ErrDecodeResponse, err)
	}

	models := make([]string, len(result.Data))
	for i, model := range result.Data {
		models[i] = model.ID
	}
	return models, nil
}

func (c *Client) GetModel(model string) (*ModelInfo, error) {
	endpoint := fmt.Sprintf("models/%s", model)
	resp, err := c.makeRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var modelInfo ModelInfo
	if err := json.NewDecoder(resp.Body).Decode(&modelInfo); err != nil {
		return nil, fmt.Errorf(resources.ErrDecodeResponse, err)
	}

	return &modelInfo, nil
}

// Chat sends a chat request to the Groq API
func (c *Client) Chat(model, message string) (*ChatResponse, error) {
	// Start timing the request
	startTime := time.Now()

	payload := map[string]interface{}{
		"model":    model,
		"messages": []map[string]string{{"role": "user", "content": message}},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf(resources.ErrEncodePayload, err)
	}

	resp, err := c.makeRequest("POST", "chat/completions", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read the entire response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	// Debug print
	// fmt.Printf("DEBUG - Raw API response: %s\n", string(bodyBytes))

	var chatResp ChatResponse
	if err := json.Unmarshal(bodyBytes, &chatResp); err != nil {
		return nil, fmt.Errorf(resources.ErrDecodeResponse, err)
	}

	// Verify that we have a valid response with content
	if len(chatResp.Choices) == 0 || chatResp.Choices[0].Message.Content == "" {
		return nil, fmt.Errorf("received empty response from API")
	}

	// Calculate elapsed time if not provided by the API
	if chatResp.Usage.CompletionTime <= 0 {
		chatResp.Usage.CompletionTime = time.Since(startTime).Seconds()
	}

	return &chatResp, nil
}

func (c *Client) makeRequest(method, endpoint string, body io.Reader) (*http.Response, error) {
	url := fmt.Sprintf("%s/%s", c.baseURL, endpoint)
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf(resources.ErrCreateRequest, err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf(resources.ErrHTTP, err)
	}

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf(resources.ErrAPI, resp.StatusCode, string(bodyBytes))
	}

	return resp, nil
}
