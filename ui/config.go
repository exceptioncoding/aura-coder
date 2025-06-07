package ui

import (
"fmt"
"strings"
"time"

"github.com/charmbracelet/bubbles/textinput"
"github.com/charmbracelet/bubbletea"
"github.com/charmbracelet/lipgloss"

"aura/config"
"aura/models"
)

var (
// Configuration screen styles
configTitleStyle = lipgloss.NewStyle().
Foreground(lipgloss.Color("#FF6B35")).
Bold(true).
Align(lipgloss.Center).
Border(lipgloss.DoubleBorder()).
BorderForeground(lipgloss.Color("#FF6B35")).
Padding(1, 4).
MarginBottom(2)

providerBoxStyle = lipgloss.NewStyle().
Foreground(lipgloss.Color("#4ECDC4")).
Border(lipgloss.RoundedBorder()).
BorderForeground(lipgloss.Color("#4ECDC4")).
Padding(1, 2).
MarginBottom(1)

selectedProviderStyle = lipgloss.NewStyle().
Foreground(lipgloss.Color("#000000")).
Background(lipgloss.Color("#4ECDC4")).
Bold(true).
Border(lipgloss.RoundedBorder()).
BorderForeground(lipgloss.Color("#4ECDC4")).
Padding(1, 2).
MarginBottom(1)

apiKeyBoxStyle = lipgloss.NewStyle().
Foreground(lipgloss.Color("#FFE66D")).
Border(lipgloss.RoundedBorder()).
BorderForeground(lipgloss.Color("#FFE66D")).
Padding(1, 2).
MarginTop(1).
MarginBottom(1)

instructionStyle = lipgloss.NewStyle().
Foreground(lipgloss.Color("#95E1D3")).
Italic(true).
MarginBottom(1)

errorConfigStyle = lipgloss.NewStyle().
Foreground(lipgloss.Color("#FF6B6B")).
Bold(true).
Border(lipgloss.RoundedBorder()).
BorderForeground(lipgloss.Color("#FF6B6B")).
Padding(1, 2).
MarginTop(1)

successConfigStyle = lipgloss.NewStyle().
Foreground(lipgloss.Color("#98FB98")).
Bold(true).
Border(lipgloss.RoundedBorder()).
BorderForeground(lipgloss.Color("#98FB98")).
Padding(1, 2).
MarginTop(1)

configSparkleFrames = []string{"⚙️", "🔧", "⚡", "🛠️", "⚙️"}
)

type LLMProvider struct {
Name        string
DisplayName string
Description string
Models      []string
APIKeyURL   string
}

var providers = []LLMProvider{
{
Name:        "anthropic",
DisplayName: "Anthropic Claude",
Description: "Advanced reasoning and coding capabilities",
Models:      []string{"claude-3-sonnet-20240229", "claude-3-haiku-20240307", "claude-3-opus-20240229"},
APIKeyURL:   "https://console.anthropic.com/",
},
{
Name:        "openai",
DisplayName: "OpenAI GPT",
Description: "Powerful language models for coding",
Models:      []string{"gpt-4", "gpt-4-turbo-preview", "gpt-3.5-turbo"},
APIKeyURL:   "https://platform.openai.com/api-keys",
},
{
Name:        "gemini",
DisplayName: "Google Gemini",
Description: "Google's advanced AI model",
Models:      []string{"gemini-pro", "gemini-pro-vision"},
APIKeyURL:   "https://makersuite.google.com/app/apikey",
},
{
Name:        "cohere",
DisplayName: "Cohere",
Description: "Enterprise-grade language AI",
Models:      []string{"command", "command-nightly"},
APIKeyURL:   "https://dashboard.cohere.ai/api-keys",
},
}

type ConfigModel struct {
selectedProvider int
apiKeyInput      textinput.Model
modelSelection   int
step             int // 0: provider, 1: api key, 2: model, 3: confirm
sparkleFrame     int
showError        string
showSuccess      bool
currentDir       string
}

func NewConfigModel(currentDir string) ConfigModel {
ti := textinput.New()
ti.Placeholder = "Enter your API key..."
ti.EchoMode = textinput.EchoPassword
ti.EchoCharacter = '•'
ti.CharLimit = 200
ti.Width = 60

return ConfigModel{
selectedProvider: 0,
apiKeyInput:      ti,
modelSelection:   0,
step:             0,
sparkleFrame:     0,
showError:        "",
showSuccess:      false,
currentDir:       currentDir,
}
}

func (m ConfigModel) Init() tea.Cmd {
return tea.Batch(
textinput.Blink,
tea.Tick(time.Millisecond*500, func(t time.Time) tea.Msg {
return models.TimerTickMsg(t)
}),
)
}

func (m ConfigModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
var cmd tea.Cmd

switch msg := msg.(type) {
case models.TimerTickMsg:
m.sparkleFrame = (m.sparkleFrame + 1) % len(configSparkleFrames)
return m, tea.Tick(time.Millisecond*500, func(t time.Time) tea.Msg {
return models.TimerTickMsg(t)
})

case tea.KeyMsg:
switch m.step {
case 0: // Provider selection
return m.handleProviderSelection(msg)
case 1: // API key input
return m.handleAPIKeyInput(msg)
case 2: // Model selection
return m.handleModelSelection(msg)
case 3: // Confirmation
return m.handleConfirmation(msg)
}
}

return m, cmd
}

func (m ConfigModel) handleProviderSelection(msg tea.KeyMsg) (ConfigModel, tea.Cmd) {
switch msg.String() {
case "up", "k":
if m.selectedProvider > 0 {
m.selectedProvider--
}
case "down", "j":
if m.selectedProvider < len(providers)-1 {
m.selectedProvider++
}
case "enter":
m.step = 1
m.apiKeyInput.Focus()
case "ctrl+c", "q":
return m, tea.Quit
}
return m, nil
}

func (m ConfigModel) handleAPIKeyInput(msg tea.KeyMsg) (ConfigModel, tea.Cmd) {
var cmd tea.Cmd

switch msg.String() {
case "ctrl+c":
return m, tea.Quit
case "esc":
m.step = 0
m.apiKeyInput.Blur()
case "enter":
apiKey := strings.TrimSpace(m.apiKeyInput.Value())
if apiKey == "" {
m.showError = "API key cannot be empty"
return m, nil
}
m.step = 2
m.apiKeyInput.Blur()
m.showError = ""
default:
m.apiKeyInput, cmd = m.apiKeyInput.Update(msg)
}
return m, cmd
}

func (m ConfigModel) handleModelSelection(msg tea.KeyMsg) (ConfigModel, tea.Cmd) {
provider := providers[m.selectedProvider]
switch msg.String() {
case "up", "k":
if m.modelSelection > 0 {
m.modelSelection--
}
case "down", "j":
if m.modelSelection < len(provider.Models)-1 {
m.modelSelection++
}
case "enter":
m.step = 3
case "esc":
m.step = 1
m.apiKeyInput.Focus()
case "ctrl+c":
return m, tea.Quit
}
return m, nil
}

func (m ConfigModel) handleConfirmation(msg tea.KeyMsg) (ConfigModel, tea.Cmd) {
switch msg.String() {
case "y", "Y":
// Save configuration
cfg := &config.Config{
LLMProvider:  providers[m.selectedProvider].Name,
APIKey:       m.apiKeyInput.Value(),
DefaultModel: providers[m.selectedProvider].Models[m.modelSelection],
TrustedDirs:  []string{},
Theme:        "default",
AutoUpdate:   true,
}

if err := cfg.Save(); err != nil {
m.showError = fmt.Sprintf("Failed to save configuration: %v", err)
m.step = 2
return m, nil
}

m.showSuccess = true
return m, tea.Tick(time.Second*2, func(_ time.Time) tea.Msg {
return models.ScreenChangeMsg{
NewScreen: models.Welcome,
Data:      m.currentDir,
}
})

case "n", "N", "esc":
m.step = 2
case "ctrl+c":
return m, tea.Quit
}
return m, nil
}

func (m ConfigModel) View() string {
var b strings.Builder

// Animated title
sparkle := configSparkleFrames[m.sparkleFrame]
titleText := fmt.Sprintf("%s Aura Configuration %s", sparkle, sparkle)
b.WriteString(configTitleStyle.Render(titleText))
b.WriteString("\n")

switch m.step {
case 0:
b.WriteString(m.renderProviderSelection())
case 1:
b.WriteString(m.renderAPIKeyInput())
case 2:
b.WriteString(m.renderModelSelection())
case 3:
b.WriteString(m.renderConfirmation())
}

if m.showError != "" {
b.WriteString("\n")
b.WriteString(errorConfigStyle.Render(fmt.Sprintf("❌ %s", m.showError)))
}

if m.showSuccess {
b.WriteString("\n")
b.WriteString(successConfigStyle.Render("✅ Configuration saved successfully! Redirecting..."))
}

return b.String()
}

func (m ConfigModel) renderProviderSelection() string {
var b strings.Builder

b.WriteString(instructionStyle.Render("🔍 Choose your preferred AI provider:"))
b.WriteString("\n\n")

for i, provider := range providers {
content := fmt.Sprintf("🤖 %s\n%s\n\n📊 Available models: %d\n🔗 Get API key: %s", 
provider.DisplayName, 
provider.Description,
len(provider.Models),
provider.APIKeyURL)

if i == m.selectedProvider {
b.WriteString("▶ ")
b.WriteString(selectedProviderStyle.Render(content))
} else {
b.WriteString("  ")
b.WriteString(providerBoxStyle.Render(content))
}
b.WriteString("\n")
}

b.WriteString("\n")
b.WriteString(instructionStyle.Render("Use ↑↓ to navigate, Enter to select, Ctrl+C to quit"))

return b.String()
}

func (m ConfigModel) renderAPIKeyInput() string {
var b strings.Builder

provider := providers[m.selectedProvider]
b.WriteString(instructionStyle.Render(fmt.Sprintf("🔑 Enter your %s API key:", provider.DisplayName)))
b.WriteString("\n\n")

b.WriteString(providerBoxStyle.Render(fmt.Sprintf("Selected: %s", provider.DisplayName)))
b.WriteString("\n\n")

keyContent := fmt.Sprintf("🔗 Get your API key from:\n%s\n\n🔐 API Key:", provider.APIKeyURL)
b.WriteString(apiKeyBoxStyle.Render(keyContent))
b.WriteString("\n")
b.WriteString(m.apiKeyInput.View())
b.WriteString("\n\n")

b.WriteString(instructionStyle.Render("Enter your API key, Esc to go back, Ctrl+C to quit"))

return b.String()
}

func (m ConfigModel) renderModelSelection() string {
var b strings.Builder

provider := providers[m.selectedProvider]
b.WriteString(instructionStyle.Render(fmt.Sprintf("🎯 Choose a %s model:", provider.DisplayName)))
b.WriteString("\n\n")

for i, model := range provider.Models {
content := fmt.Sprintf("🧠 %s", model)
if i == m.modelSelection {
b.WriteString("▶ ")
b.WriteString(selectedProviderStyle.Render(content))
} else {
b.WriteString("  ")
b.WriteString(providerBoxStyle.Render(content))
}
b.WriteString("\n")
}

b.WriteString("\n")
b.WriteString(instructionStyle.Render("Use ↑↓ to navigate, Enter to continue, Esc to go back"))

return b.String()
}

func (m ConfigModel) renderConfirmation() string {
var b strings.Builder

provider := providers[m.selectedProvider]
selectedModel := provider.Models[m.modelSelection]

b.WriteString(instructionStyle.Render("✅ Confirm your configuration:"))
b.WriteString("\n\n")

configSummary := fmt.Sprintf(`🤖 Provider: %s
🧠 Model: %s
🔑 API Key: %s
📁 Directory: %s`,
provider.DisplayName,
selectedModel,
strings.Repeat("•", len(m.apiKeyInput.Value())),
m.currentDir)

b.WriteString(successConfigStyle.Render(configSummary))
b.WriteString("\n\n")

b.WriteString(instructionStyle.Render("Save this configuration? [y/N]"))

return b.String()
}
