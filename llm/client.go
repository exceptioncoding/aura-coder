package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"aura/config"
)

type Client struct {
	config *config.Config
	client *http.Client
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
}

type ChatResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

// Anthropic specific types
type AnthropicRequest struct {
	Model     string    `json:"model"`
	MaxTokens int       `json:"max_tokens"`
	Messages  []Message `json:"messages"`
	System    string    `json:"system,omitempty"`
}

type AnthropicResponse struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Role    string `json:"role"`
	Content []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"content"`
	Model        string `json:"model"`
	StopReason   string `json:"stop_reason"`
	StopSequence string `json:"stop_sequence"`
	Usage        struct {
		InputTokens  int `json:"input_tokens"`
		OutputTokens int `json:"output_tokens"`
	} `json:"usage"`
}

func NewClient(cfg *config.Config) *Client {
	return &Client{
		config: cfg,
		client: &http.Client{
			Timeout: time.Second * 60,
		},
	}
}

func (c *Client) Chat(messages []Message) (string, error) {
	switch c.config.LLMProvider {
	case "openai":
		return c.chatOpenAI(messages)
	case "anthropic":
		return c.chatAnthropic(messages)
	case "gemini":
		return c.chatGemini(messages)
	case "cohere":
		return c.chatCohere(messages)
	default:
		return "", fmt.Errorf("unsupported LLM provider: %s", c.config.LLMProvider)
	}
}

func (c *Client) chatOpenAI(messages []Message) (string, error) {
	req := ChatRequest{
		Model:    c.config.DefaultModel,
		Messages: messages,
		Stream:   false,
	}

	reqBody, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(reqBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
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

func (c *Client) chatAnthropic(messages []Message) (string, error) {
	// Separate system message if present
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
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", "https://api.anthropic.com/v1/messages", bytes.NewBuffer(reqBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
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

func (c *Client) chatGemini(_ []Message) (string, error) {
	// Placeholder for Gemini implementation
	return "", fmt.Errorf("gemini integration coming soon! Please use OpenAI or Anthropic for now")
}

func (c *Client) chatCohere(_ []Message) (string, error) {
	// Placeholder for Cohere implementation
	return "", fmt.Errorf("cohere integration coming soon! Please use OpenAI or Anthropic for now")
}

// Helper function to create a coding assistant system prompt
func (c *Client) CreateCodingPrompt(userInput, currentDir string) []Message {
	systemPrompt := fmt.Sprintf(`You are Aura, an advanced AI coding assistant. You are currently working in the directory: %s

Your capabilities include:
- Analyzing code and providing insights
- Writing and refactoring code
- Debugging and troubleshooting
- Explaining programming concepts
- Suggesting best practices
- Helping with git operations
- Running commands (with user approval)

Guidelines:
- Be concise but comprehensive
- Provide code examples when helpful
- Ask for clarification when needed
- Suggest specific actionable steps
- Focus on best practices and clean code
- Consider security implications

Current working directory: %s`, currentDir, currentDir)

	return []Message{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: userInput},
	}
}
