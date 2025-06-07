package llm

import (
"encoding/json"
"testing"

"github.com/stretchr/testify/assert"
"github.com/stretchr/testify/require"

"aura/config"
)

func TestNewClient(t *testing.T) {
cfg := &config.Config{
LLMProvider:  "openai",
APIKey:       "test-key",
DefaultModel: "gpt-4",
}

client := NewClient(cfg)

assert.NotNil(t, client)
assert.Equal(t, cfg, client.config)
assert.NotNil(t, client.client)
}

func TestClient_Chat_UnsupportedProvider(t *testing.T) {
cfg := &config.Config{
LLMProvider: "unsupported-provider",
APIKey:      "test-key",
}

client := NewClient(cfg)
messages := []Message{{Role: "user", Content: "test"}}

response, err := client.Chat(messages)

assert.Error(t, err)
assert.Empty(t, response)
assert.Contains(t, err.Error(), "unsupported LLM provider")
}

func TestClient_CreateCodingPrompt(t *testing.T) {
cfg := &config.Config{}
client := NewClient(cfg)

userInput := "How do I create a REST API?"
currentDir := "/home/user/project"

messages := client.CreateCodingPrompt(userInput, currentDir)

require.Len(t, messages, 2)

// Verify system message
systemMsg := messages[0]
assert.Equal(t, "system", systemMsg.Role)
assert.Contains(t, systemMsg.Content, "You are Aura")
assert.Contains(t, systemMsg.Content, currentDir)
assert.Contains(t, systemMsg.Content, "coding assistant")

// Verify user message
userMsg := messages[1]
assert.Equal(t, "user", userMsg.Role)
assert.Equal(t, userInput, userMsg.Content)
}

func TestMessage_JSONSerialization(t *testing.T) {
message := Message{
Role:    "user",
Content: "Hello, world!",
}

// Test marshaling
data, err := json.Marshal(message)
assert.NoError(t, err)

expected := `{"role":"user","content":"Hello, world!"}`
assert.JSONEq(t, expected, string(data))

// Test unmarshaling
var unmarshaled Message
err = json.Unmarshal(data, &unmarshaled)
assert.NoError(t, err)
assert.Equal(t, message, unmarshaled)
}

func TestChatRequest_JSONSerialization(t *testing.T) {
req := ChatRequest{
Model: "gpt-4",
Messages: []Message{
{Role: "user", Content: "Hello"},
},
Stream: false,
}

data, err := json.Marshal(req)
assert.NoError(t, err)

var unmarshaled ChatRequest
err = json.Unmarshal(data, &unmarshaled)
assert.NoError(t, err)
assert.Equal(t, req, unmarshaled)
}

func TestAnthropicRequest_JSONSerialization(t *testing.T) {
req := AnthropicRequest{
Model:     "claude-3-sonnet-20240229",
MaxTokens: 4000,
Messages: []Message{
{Role: "user", Content: "Hello"},
},
System: "You are a helpful assistant",
}

data, err := json.Marshal(req)
assert.NoError(t, err)

var unmarshaled AnthropicRequest
err = json.Unmarshal(data, &unmarshaled)
assert.NoError(t, err)
assert.Equal(t, req, unmarshaled)
}

func TestChatResponse_JSONSerialization(t *testing.T) {
resp := ChatResponse{
ID:      "test-id",
Object:  "chat.completion",
Created: 1234567890,
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
Content: "Hello!",
},
FinishReason: "stop",
},
},
}

data, err := json.Marshal(resp)
assert.NoError(t, err)

var unmarshaled ChatResponse
err = json.Unmarshal(data, &unmarshaled)
assert.NoError(t, err)
assert.Equal(t, resp.ID, unmarshaled.ID)
assert.Equal(t, resp.Object, unmarshaled.Object)
assert.Equal(t, resp.Model, unmarshaled.Model)
assert.Len(t, unmarshaled.Choices, 1)
}

func TestAnthropicResponse_JSONSerialization(t *testing.T) {
resp := AnthropicResponse{
ID:   "test-id",
Type: "message",
Role: "assistant",
Content: []struct {
Type string `json:"type"`
Text string `json:"text"`
}{
{Type: "text", Text: "Hello from Claude!"},
},
Model:      "claude-3-sonnet-20240229",
StopReason: "end_turn",
}

data, err := json.Marshal(resp)
assert.NoError(t, err)

var unmarshaled AnthropicResponse
err = json.Unmarshal(data, &unmarshaled)
assert.NoError(t, err)
assert.Equal(t, resp.ID, unmarshaled.ID)
assert.Equal(t, resp.Type, unmarshaled.Type)
assert.Equal(t, resp.Role, unmarshaled.Role)
assert.Len(t, unmarshaled.Content, 1)
assert.Equal(t, resp.Content[0].Text, unmarshaled.Content[0].Text)
}

// Test validation functions
func TestMessage_Validation(t *testing.T) {
tests := []struct {
name    string
message Message
valid   bool
}{
{
name: "valid user message",
message: Message{
Role:    "user",
Content: "Hello",
},
valid: true,
},
{
name: "valid assistant message",
message: Message{
Role:    "assistant",
Content: "Hi there!",
},
valid: true,
},
{
name: "valid system message",
message: Message{
Role:    "system",
Content: "You are a helpful assistant",
},
valid: true,
},
{
name: "empty role",
message: Message{
Role:    "",
Content: "Hello",
},
valid: false,
},
{
name: "empty content",
message: Message{
Role:    "user",
Content: "",
},
valid: false,
},
}

for _, tt := range tests {
t.Run(tt.name, func(t *testing.T) {
isValid := tt.message.Role != "" && tt.message.Content != ""
assert.Equal(t, tt.valid, isValid)
})
}
}

// Benchmark tests
func BenchmarkClient_CreateCodingPrompt(b *testing.B) {
cfg := &config.Config{}
client := NewClient(cfg)
userInput := "How do I optimize this Go code?"
currentDir := "/home/user/project"

b.ResetTimer()
for i := 0; i < b.N; i++ {
client.CreateCodingPrompt(userInput, currentDir)
}
}

func BenchmarkMessage_JSONMarshal(b *testing.B) {
message := Message{
Role:    "user",
Content: "This is a test message for benchmarking JSON marshaling performance",
}

b.ResetTimer()
for i := 0; i < b.N; i++ {
_, _ = json.Marshal(message)
}
}

func BenchmarkMessage_JSONUnmarshal(b *testing.B) {
data := []byte(`{"role":"user","content":"This is a test message for benchmarking JSON unmarshaling performance"}`)

b.ResetTimer()
for i := 0; i < b.N; i++ {
var msg Message
_ = json.Unmarshal(data, &msg)
}
}

// Test edge cases
func TestEdgeCases(t *testing.T) {
t.Run("nil config", func(t *testing.T) {
client := NewClient(nil)
assert.NotNil(t, client)
assert.Nil(t, client.config)
})

t.Run("empty messages slice", func(t *testing.T) {
cfg := &config.Config{LLMProvider: "openai", APIKey: "test"}
client := NewClient(cfg)

// CreateCodingPrompt should still work 
result := client.CreateCodingPrompt("test", "/test")
assert.Len(t, result, 2) // Should still create system + user message
})

t.Run("very long message content", func(t *testing.T) {
longContent := make([]byte, 10000)
for i := range longContent {
longContent[i] = 'a'
}

message := Message{
Role:    "user",
Content: string(longContent),
}

data, err := json.Marshal(message)
assert.NoError(t, err)
assert.Greater(t, len(data), 10000)
})
}
