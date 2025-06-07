package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"aura/models"
)

var (
	// Welcome screen styles
	welcomeTitleStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#9D4EDD")).
				Bold(true).
				Align(lipgloss.Center).
				Border(lipgloss.DoubleBorder()).
				BorderForeground(lipgloss.Color("#9D4EDD")).
				Padding(1, 4).
				MarginBottom(2)

	dirStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#06FFA5")).
			Bold(true).
			Background(lipgloss.Color("#001122")).
			Padding(0, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#06FFA5"))

	featureStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFE66D")).
			MarginLeft(2)

	tipStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF6B9D")).
			Bold(true).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#FF6B9D")).
			Padding(1, 2).
			MarginTop(2)

	continueStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#4ECDC4")).
			Bold(true).
			Align(lipgloss.Center).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#4ECDC4")).
			Padding(0, 2).
			Blink(true)

	sparkleFrames = []string{"✨", "🌟", "💫", "⭐", "✨"}
)

type WelcomeModel struct {
	currentDir   string
	sparkleFrame int
	fadeIn       bool
	fadeOpacity  float64
	showContinue bool
}

func NewWelcomeModel(currentDir string) WelcomeModel {
	return WelcomeModel{
		currentDir:   currentDir,
		sparkleFrame: 0,
		fadeIn:       true,
		fadeOpacity:  0.0,
		showContinue: false,
	}
}

func (m WelcomeModel) Init() tea.Cmd {
	return tea.Batch(
		tea.Tick(time.Millisecond*300, func(t time.Time) tea.Msg {
			return models.TimerTickMsg(t)
		}),
	)
}

func (m WelcomeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case models.TimerTickMsg:
		// Handle sparkle animation
		m.sparkleFrame = (m.sparkleFrame + 1) % len(sparkleFrames)

		// Handle fade-in effect
		if m.fadeIn && m.fadeOpacity < 1.0 {
			m.fadeOpacity += 0.1
			if m.fadeOpacity >= 1.0 {
				m.fadeOpacity = 1.0
				m.fadeIn = false
				m.showContinue = true
			}
		}

		return m, tea.Tick(time.Millisecond*300, func(t time.Time) tea.Msg {
			return models.TimerTickMsg(t)
		})

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, key.NewBinding(key.WithKeys("q", "ctrl+c"))):
			return m, tea.Quit

		case key.Matches(msg, key.NewBinding(key.WithKeys("enter", " "))):
			return m, func() tea.Msg {
				return models.ScreenChangeMsg{
					NewScreen: models.Main,
					Data:      m.currentDir,
				}
			}
		}
	}

	return m, nil
}

func (m WelcomeModel) View() string {
	var b strings.Builder

	// Animated title with sparkles
	sparkle := sparkleFrames[m.sparkleFrame]
	titleText := fmt.Sprintf("%s Welcome to Aura %s", sparkle, sparkle)
	b.WriteString(welcomeTitleStyle.Render(titleText))
	b.WriteString("\n")

	// Current directory with nice styling
	b.WriteString("Current working directory:\n")
	b.WriteString(dirStyle.Render(fmt.Sprintf("📁 %s", m.currentDir)))
	b.WriteString("\n\n")

	// Features list with icons
	features := []string{
		"🤖 Type natural language commands to interact with your codebase",
		"⚡ Use /help to see available slash commands",
		"📊 Use /status to check your current configuration",
		"🔒 All file changes and command executions require confirmation",
		"🔄 Git integration for seamless version control",
		"💡 AI-powered code analysis and generation",
	}

	b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#FFE66D")).Bold(true).Render("🚀 Getting Started:"))
	b.WriteString("\n")

	for _, feature := range features {
		b.WriteString(featureStyle.Render(fmt.Sprintf("  %s", feature)))
		b.WriteString("\n")
	}

	// Pro tip
	tipText := "💡 Pro Tip: Start with '/help' to explore all available commands, or just ask Aura to help you with your coding tasks!"
	b.WriteString("\n")
	b.WriteString(tipStyle.Render(tipText))

	// Continue prompt with animation
	if m.showContinue {
		b.WriteString("\n\n")
		b.WriteString(continueStyle.Render("Press Enter to start coding with Aura! 🚀"))
	}

	return b.String()
}
