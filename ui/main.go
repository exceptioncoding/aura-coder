package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"aura/config"
	"aura/llm"
	"aura/models"
)

var (
	// Main interface styles
	headerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF6B35")).
			Background(lipgloss.Color("#1A1A2E")).
			Bold(true).
			Padding(0, 2).
			MarginBottom(1)

	dirStyle2 = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#0FB9B1")).
			Bold(true)

	separatorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#495057"))

	userMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#4ECDC4")).
				Bold(true)

	auraMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFE66D")).
				MarginLeft(2)

	systemMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#95E1D3")).
				Italic(true).
				MarginLeft(2)

	errorMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FF6B6B")).
				Bold(true).
				MarginLeft(2)

	confirmationStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFE66D")).
				Background(lipgloss.Color("#2D3748")).
				Bold(true).
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("#FFE66D")).
				Padding(1, 2).
				MarginTop(1)

	promptStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#4ECDC4")).
			Bold(true)

	helpBoxStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#A8E6CF")).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#A8E6CF")).
			Padding(1, 2).
			MarginTop(1).
			MarginBottom(1)

	thinkingFrames = []string{"🤔", "💭", "🧠", "⚡", "💡"}
	loadingFrames  = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
)

type MainModel struct {
	currentDir       string
	textInput        textinput.Model
	chatHistory      []ChatMessage
	pendingCommand   *models.PendingCommand
	showConfirmation bool
	thinkingFrame    int
	loadingFrame     int
	isProcessing     bool
	terminalWidth    int
	terminalHeight   int
}

type ChatMessage struct {
	Type    string // "user", "aura", "system", "error"
	Content string
	Time    time.Time
}

func NewMainModel(currentDir string) MainModel {
	ti := textinput.New()
	ti.Placeholder = "Type your message..."
	ti.Focus()
	ti.CharLimit = 500
	ti.Width = 60 // Safe initial width

	return MainModel{
		currentDir:       currentDir,
		textInput:        ti,
		chatHistory:      []ChatMessage{},
		pendingCommand:   nil,
		showConfirmation: false,
		thinkingFrame:    0,
		loadingFrame:     0,
		isProcessing:     false,
		terminalWidth:    80,
		terminalHeight:   24,
	}
}

func (m MainModel) Init() tea.Cmd {
	return tea.Batch(
		textinput.Blink,
		tea.Tick(time.Millisecond*150, func(t time.Time) tea.Msg {
			return models.TimerTickMsg(t)
		}),
	)
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case models.TimerTickMsg:
		m.thinkingFrame = (m.thinkingFrame + 1) % len(thinkingFrames)
		m.loadingFrame = (m.loadingFrame + 1) % len(loadingFrames)
		return m, tea.Tick(time.Millisecond*150, func(t time.Time) tea.Msg {
			return models.TimerTickMsg(t)
		})

	case tea.WindowSizeMsg:
		m.terminalWidth = msg.Width
		m.terminalHeight = msg.Height
		// Ensure minimum width to prevent panic
		if msg.Width > 10 {
			m.textInput.Width = msg.Width - 4
		} else {
			m.textInput.Width = 6
		}

	case tea.KeyMsg:
		if m.showConfirmation {
			return m.handleConfirmation(msg)
		}

		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			return m.handleInput()
		case "ctrl+l":
			m.chatHistory = []ChatMessage{}
			return m, nil
		default:
			m.textInput, cmd = m.textInput.Update(msg)
			return m, cmd
		}

	case models.ConfirmationMsg:
		m.pendingCommand = &models.PendingCommand{
			Description: msg.Prompt,
			Command:     msg.Command,
		}
		m.showConfirmation = true
		return m, nil

	case models.StatusMsg:
		m.isProcessing = false
		m.chatHistory = append(m.chatHistory, ChatMessage{
			Type:    "aura",
			Content: msg.Message,
			Time:    time.Now(),
		})
		return m, nil
	}

	return m, nil
}

func (m MainModel) handleInput() (MainModel, tea.Cmd) {
	input := strings.TrimSpace(m.textInput.Value())
	if input == "" {
		return m, nil
	}

	// Add user message
	m.chatHistory = append(m.chatHistory, ChatMessage{
		Type:    "user",
		Content: input,
		Time:    time.Now(),
	})

	m.textInput.SetValue("")

	// Handle special commands without slash
	if strings.ToLower(input) == "exit" {
		return m, tea.Quit
	}

	if strings.HasPrefix(input, "/") {
		return m.handleSlashCommand(input)
	}

	// Handle natural language input with LLM
	m.isProcessing = true
	m.chatHistory = append(m.chatHistory, ChatMessage{
		Type:    "system",
		Content: "Processing your request...",
		Time:    time.Now(),
	})

	// Call LLM
	return m, func() tea.Msg {
		cfg, err := config.LoadConfig()
		if err != nil || cfg.APIKey == "" {
			return models.StatusMsg{Message: "❌ Configuration error. Please run the configuration setup by restarting Aura."}
		}

		client := llm.NewClient(cfg)
		messages := client.CreateCodingPrompt(input, m.currentDir)

		response, err := client.Chat(messages)
		if err != nil {
			return models.StatusMsg{Message: fmt.Sprintf("❌ AI Error: %v", err)}
		}

		return models.StatusMsg{Message: response}
	}

}

func (m MainModel) handleSlashCommand(command string) (MainModel, tea.Cmd) {
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return m, nil
	}

	switch parts[0] {
	case "/help":
		m.chatHistory = append(m.chatHistory, ChatMessage{
			Type:    "aura",
			Content: m.getHelpText(),
			Time:    time.Now(),
		})
	case "/status":
		m.chatHistory = append(m.chatHistory, ChatMessage{
			Type:    "aura",
			Content: m.getStatusText(),
			Time:    time.Now(),
		})
	case "/clear":
		m.chatHistory = []ChatMessage{}
	case "/exit":
		return m, tea.Quit
	case "/edit":
		if len(parts) > 1 {
			m.chatHistory = append(m.chatHistory, ChatMessage{
				Type:    "aura",
				Content: fmt.Sprintf("📝 Opening %s for editing...", parts[1]),
				Time:    time.Now(),
			})
		} else {
			m.chatHistory = append(m.chatHistory, ChatMessage{
				Type:    "error",
				Content: "❌ Please specify a file to edit. Usage: /edit <filename>",
				Time:    time.Now(),
			})
		}
	case "/run":
		if len(parts) > 1 {
			cmd := strings.Join(parts[1:], " ")
			return m, func() tea.Msg {
				return models.ConfirmationMsg{
					Prompt:  fmt.Sprintf("🚀 Execute command: %s", cmd),
					Command: cmd,
				}
			}
		} else {
			m.chatHistory = append(m.chatHistory, ChatMessage{
				Type:    "error",
				Content: "❌ Please specify a command to run. Usage: /run <command>",
				Time:    time.Now(),
			})
		}
	case "/test":
		return m, func() tea.Msg {
			return models.ConfirmationMsg{
				Prompt:  "🧪 Run project tests",
				Command: "go test ./...",
			}
		}
	default:
		m.chatHistory = append(m.chatHistory, ChatMessage{
			Type:    "error",
			Content: fmt.Sprintf("❌ Unknown command: %s. Type /help for available commands.", parts[0]),
			Time:    time.Now(),
		})
	}

	return m, nil
}

func (m MainModel) handleConfirmation(msg tea.KeyMsg) (MainModel, tea.Cmd) {
	switch msg.String() {
	case "y", "Y":
		m.chatHistory = append(m.chatHistory, ChatMessage{
			Type:    "system",
			Content: fmt.Sprintf("✅ Executing: %s", m.pendingCommand.Command),
			Time:    time.Now(),
		})
		m.showConfirmation = false
		m.pendingCommand = nil
	case "n", "N", "escape":
		m.chatHistory = append(m.chatHistory, ChatMessage{
			Type:    "system",
			Content: "❌ Command cancelled by user",
			Time:    time.Now(),
		})
		m.showConfirmation = false
		m.pendingCommand = nil
	}
	return m, nil
}

func (m MainModel) getHelpText() string {
	return `🔧 Available Commands:

💬 Chat Commands:
  /help       - Show this help message
  /status     - Show current status and configuration
  /clear      - Clear chat history
  /exit       - Exit Aura (or type "exit")

📁 File Operations:
  /edit <file> - Open file for AI-assisted editing

⚡ System Commands:
  /run <cmd>   - Execute a shell command (with confirmation)
  /test        - Run project tests
  /commit      - Interactive Git commit

💡 Tips:
  • Type natural language requests to interact with your codebase
  • Use Ctrl+L to quickly clear the screen
  • All operations require confirmation for safety
  • Configure your API key in ~/.aura.conf for AI features`
}

func (m MainModel) getStatusText() string {
	return fmt.Sprintf(`📊 System Status:

📂 Current Directory: %s
⚙️  Configuration: ~/.aura.conf
🤖 LLM Provider: Not configured (set API key in config)
🔧 Terminal Size: %dx%d
💾 Chat History: %d messages

🚀 Ready to assist with your coding tasks!`,
		m.currentDir, m.terminalWidth, m.terminalHeight, len(m.chatHistory))
}

func (m MainModel) View() string {
	var b strings.Builder

	// Animated header
	headerText := fmt.Sprintf("🌟 Aura - AI Coding Agent %s", dirStyle2.Render(fmt.Sprintf("| %s", m.currentDir)))
	b.WriteString(headerStyle.Render(headerText))
	b.WriteString("\n")
	// Ensure separator doesn't exceed reasonable width
	separatorWidth := m.terminalWidth
	if separatorWidth > 120 {
		separatorWidth = 120
	}
	if separatorWidth < 20 {
		separatorWidth = 20
	}
	b.WriteString(separatorStyle.Render(strings.Repeat("─", separatorWidth)))
	b.WriteString("\n\n")

	// Chat history with better styling
	maxMessages := (m.terminalHeight - 8) // Reserve space for header and input
	startIdx := 0
	if len(m.chatHistory) > maxMessages {
		startIdx = len(m.chatHistory) - maxMessages
	}

	for i := startIdx; i < len(m.chatHistory); i++ {
		msg := m.chatHistory[i]
		timeStr := msg.Time.Format("15:04")

		switch msg.Type {
		case "user":
			b.WriteString(userMessageStyle.Render(fmt.Sprintf("👤 [%s] %s", timeStr, msg.Content)))
		case "aura":
			b.WriteString(auraMessageStyle.Render(fmt.Sprintf("🤖 [%s] %s", timeStr, msg.Content)))
		case "system":
			if m.isProcessing && i == len(m.chatHistory)-1 {
				thinkingIcon := thinkingFrames[m.thinkingFrame]
				b.WriteString(systemMessageStyle.Render(fmt.Sprintf("%s [%s] %s", thinkingIcon, timeStr, msg.Content)))
			} else {
				b.WriteString(systemMessageStyle.Render(fmt.Sprintf("ℹ️  [%s] %s", timeStr, msg.Content)))
			}
		case "error":
			b.WriteString(errorMessageStyle.Render(fmt.Sprintf("⚠️  [%s] %s", timeStr, msg.Content)))
		}
		b.WriteString("\n")
	}

	if len(m.chatHistory) > 0 {
		b.WriteString("\n")
	}

	// Confirmation prompt with nice styling
	if m.showConfirmation && m.pendingCommand != nil {
		confirmText := fmt.Sprintf("🔐 %s\n\n💡 This action requires your confirmation.\nProceed? [y/N]", m.pendingCommand.Description)
		b.WriteString(confirmationStyle.Render(confirmText))
		return b.String()
	}

	// Input prompt with loading animation
	promptText := "> "
	if m.isProcessing {
		loadingIcon := loadingFrames[m.loadingFrame]
		promptText = fmt.Sprintf("%s ", loadingIcon)
	}

	b.WriteString(promptStyle.Render(promptText))
	b.WriteString(m.textInput.View())

	// Show helpful hints for new users
	if len(m.chatHistory) == 0 {
		b.WriteString("\n\n")
		hints := "💡 Quick Start: Type '/help' for commands, or just ask me about your code!"
		b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#888888")).Render(hints))
	}

	return b.String()
}
