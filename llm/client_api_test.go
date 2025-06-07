package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"aura/config"
)

// Test actual API interactions with mocked HTTP servers

func TestClient_ChatOpenAI_Success(t *testing.T) {
	// Create mock OpenAI API server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request headers
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "Bearer test-api-key", r.Header.Get("Authorization"))
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/v1/chat/completions", r.URL.Path)

		// Verify request body
		var req ChatRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		require.NoError(t, err)
		assert.Equal(t, "gpt-4", req.Model)
		assert.Len(t, req.Messages, 1)
		assert.Equal(t, "user", req.Messages[0].Role)
		assert.Equal(t, "Hello, AI!", req.Messages[0].Content)

		// Send mock response
		response := ChatResponse{
			ID:      "chatcmpl-123",
			Object:  "chat.completion",
			Created: 1677652288,
			Model:   "gpt-4",
			Choices: []struct {
				Index   int `json:"index"`
				Message struct {
					Role    string `json:"role"`
					Content string `json:"content"`
				} `json:"message"`
				FinishReason string `json:"finish_reason"`
			}{
				{
					Index: 0,
					Message: struct {
						Role    string `json:"role"`
						Content string `json:"content"`
					}{
						Role:    "assistant",
						Content: "Hello! How can I help you today?",
					},
					FinishReason: "stop",
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create client with mock server URL
	cfg := &config.Config{
		LLMProvider:  "openai",
		APIKey:       "test-api-key",
		DefaultModel: "gpt-4",
	}
	client := NewClient(cfg)

	// We use a test-specific method that accepts a custom URL
	response, err := client.testChatOpenAI([]Message{
		{Role: "user", Content: "Hello, AI!"},
	}, server.URL+"/v1/chat/completions")

	assert.NoError(t, err)
	assert.Equal(t, "Hello! How can I help you today?", response)
}

func TestClient_ChatAnthropic_Success(t *testing.T) {
	// Create mock Anthropic API server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request headers
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "test-api-key", r.Header.Get("x-api-key"))
		assert.Equal(t, "2023-06-01", r.Header.Get("anthropic-version"))
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/v1/messages", r.URL.Path)

		// Verify request body
		var req AnthropicRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		require.NoError(t, err)
		assert.Equal(t, "claude-3-sonnet-20240229", req.Model)
		assert.Equal(t, 4000, req.MaxTokens)
		assert.Len(t, req.Messages, 1)
		assert.Equal(t, "user", req.Messages[0].Role)

		// Send mock response
		response := AnthropicResponse{
			ID:   "msg_123",
			Type: "message",
			Role: "assistant",
			Content: []struct {
				Type string `json:"type"`
				Text string `json:"text"`
			}{
				{Type: "text", Text: "Hello! I'm Claude, how can I assist you?"},
			},
			Model:      "claude-3-sonnet-20240229",
			StopReason: "end_turn",
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create client
	cfg := &config.Config{
		LLMProvider:  "anthropic",
		APIKey:       "test-api-key",
		DefaultModel: "claude-3-sonnet-20240229",
	}
	client := NewClient(cfg)

	// Test with mock server
	response, err := client.testChatAnthropic([]Message{
		{Role: "user", Content: "Hello, Claude!"},
	}, server.URL+"/v1/messages")

	assert.NoError(t, err)
	assert.Equal(t, "Hello! I'm Claude, how can I assist you?", response)
}

func TestClient_ChatOpenAI_APIError(t *testing.T) {
	// Create mock server that returns an error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error": {"message": "Invalid API key"}}`))
	}))
	defer server.Close()

	cfg := &config.Config{
		LLMProvider:  "openai",
		APIKey:       "invalid-key",
		DefaultModel: "gpt-4",
	}
	client := NewClient(cfg)

	response, err := client.testChatOpenAI([]Message{
		{Role: "user", Content: "Hello"},
	}, server.URL+"/v1/chat/completions")

	assert.Error(t, err)
	assert.Empty(t, response)
	assert.Contains(t, err.Error(), "API request failed with status 401")
}

func TestClient_ChatAnthropic_NetworkError(t *testing.T) {
	// Test with invalid URL to simulate network error
	cfg := &config.Config{
		LLMProvider:  "anthropic",
		APIKey:       "test-key",
		DefaultModel: "claude-3-sonnet-20240229",
	}
	client := NewClient(cfg)

	response, err := client.testChatAnthropic([]Message{
		{Role: "user", Content: "Hello"},
	}, "http://invalid-url-that-does-not-exist.com/v1/messages")

	assert.Error(t, err)
	assert.Empty(t, response)
	assert.Contains(t, err.Error(), "failed to send request")
}

func TestClient_ChatOpenAI_EmptyResponse(t *testing.T) {
	// Create mock server that returns empty choices
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		response := ChatResponse{
			ID:      "chatcmpl-123",
			Object:  "chat.completion",
			Created: 1677652288,
			Model:   "gpt-4",
			Choices: []struct {
				Index   int `json:"index"`
				Message struct {
					Role    string `json:"role"`
					Content string `json:"content"`
				} `json:"message"`
				FinishReason string `json:"finish_reason"`
			}{}, // Empty choices
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	cfg := &config.Config{
		LLMProvider:  "openai",
		APIKey:       "test-key",
		DefaultModel: "gpt-4",
	}
	client := NewClient(cfg)

	response, err := client.testChatOpenAI([]Message{
		{Role: "user", Content: "Hello"},
	}, server.URL+"/v1/chat/completions")

	assert.Error(t, err)
	assert.Empty(t, response)
	assert.Contains(t, err.Error(), "no choices in response")
}

func TestClient_Chat_Integration_WithMocking(t *testing.T) {
	tests := []struct {
		name          string
		provider      string
		expectedError string
		shouldSucceed bool
	}{
		{
			name:          "unsupported provider",
			provider:      "invalid-provider",
			expectedError: "unsupported LLM provider",
			shouldSucceed: false,
		},
		{
			name:          "gemini not implemented",
			provider:      "gemini",
			expectedError: "gemini integration coming soon",
			shouldSucceed: false,
		},
		{
			name:          "cohere not implemented",
			provider:      "cohere",
			expectedError: "cohere integration coming soon",
			shouldSucceed: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{
				LLMProvider: tt.provider,
				APIKey:      "test-key",
			}
			client := NewClient(cfg)

			response, err := client.Chat([]Message{
				{Role: "user", Content: "Hello"},
			})

			if tt.shouldSucceed {
				assert.NoError(t, err)
				assert.NotEmpty(t, response)
			} else {
				assert.Error(t, err)
				assert.Empty(t, response)
				if tt.expectedError != "" {
					assert.Contains(t, err.Error(), tt.expectedError)
				}
			}
		})
	}
}

// Test helper methods for API testing with custom URLs
func (c *Client) testChatOpenAI(messages []Message, url string) (string, error) {
	req := ChatRequest{
		Model:    c.config.DefaultModel,
		Messages: messages,
		Stream:   false,
	}

	reqBody, err := json.Marshal(req)
	if err != nil {
		return "", err
	}

	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return "", err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.config.APIKey)

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var chatResp ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("no choices in response")
	}

	return chatResp.Choices[0].Message.Content, nil
}

func (c *Client) testChatAnthropic(messages []Message, url string) (string, error) {
	var systemMsg string
	var userMessages []Message
	for _, msg := range messages {
		if msg.Role == "system" {
			systemMsg = msg.Content
		} else {
			userMessages = append(userMessages, msg)
		}
	}

	req := AnthropicRequest{
		Model:     c.config.DefaultModel,
		MaxTokens: 4000,
		Messages:  userMessages,
		System:    systemMsg,
	}

	reqBody, err := json.Marshal(req)
	if err != nil {
		return "", err
	}

	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return "", err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-api-key", c.config.APIKey)
	httpReq.Header.Set("anthropic-version", "2023-06-01")

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var anthropicResp AnthropicResponse
	if err := json.NewDecoder(resp.Body).Decode(&anthropicResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if len(anthropicResp.Content) == 0 {
		return "", fmt.Errorf("no content in response")
	}

	return anthropicResp.Content[0].Text, nil
}

// Benchmark API serialization performance
func BenchmarkAPI_OpenAI_RequestSerialization(b *testing.B) {
	req := ChatRequest{
		Model: "gpt-4",
		Messages: []Message{
			{Role: "system", Content: "You are a helpful assistant"},
			{Role: "user", Content: "Write a function to reverse a string"},
		},
		Stream: false,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := json.Marshal(req)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkAPI_Anthropic_RequestSerialization(b *testing.B) {
	req := AnthropicRequest{
		Model:     "claude-3-sonnet-20240229",
		MaxTokens: 4000,
		Messages: []Message{
			{Role: "user", Content: "Write a function to reverse a string"},
		},
		System: "You are a helpful assistant",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := json.Marshal(req)
		if err != nil {
			b.Fatal(err)
		}
	}
}
