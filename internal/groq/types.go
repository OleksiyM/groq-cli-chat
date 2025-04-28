package groq

// ChatResponse represents the structure of a chat completion response
type ChatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Usage struct {
		TotalTokens    int     `json:"total_tokens"`
		CompletionTime float64 `json:"completion_time"`
	} `json:"usage"`
}

// ModelInfo represents the structure of a model retrieval response
type ModelInfo struct {
	ID            string      `json:"id"`
	Object        string      `json:"object"`
	Created       int64       `json:"created"`
	OwnedBy       string      `json:"owned_by"`
	Active        bool        `json:"active"`
	ContextWindow int         `json:"context_window"`
	PublicApps    interface{} `json:"public_apps"`
}