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
	// Security prompt styles
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFD700")).
			Bold(true).
			Align(lipgloss.Center).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#FFD700")).
			Padding(0, 2).
			MarginBottom(1)

	directoryStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#87CEEB")).
			Bold(true).
			Background(lipgloss.Color("#1E1E1E")).
			Padding(0, 1).
			MarginLeft(2)

	warningStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF6B6B")).
			Bold(true).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#FF6B6B")).
			Padding(1, 2).
			MarginTop(1).
			MarginBottom(1)


	selectedChoiceStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#000000")).
				Background(lipgloss.Color("#98FB98")).
				Bold(true).
				Padding(0, 1).
				MarginLeft(1)

	unselectedChoiceStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#98FB98")).
				Bold(true).
				MarginLeft(2)

	pulseColors = []string{"#FFD700", "#FFA500", "#FF8C00", "#FFD700"}
)

type SecurityPromptModel struct {
	directory      string
	selectedChoice int
	pulseFrame     int
	showAnimation  bool
}

func NewSecurityPromptModel(directory string) SecurityPromptModel {
	return SecurityPromptModel{
		directory:      directory,
		selectedChoice: 0,
		pulseFrame:     0,
		showAnimation:  true,
	}
}

func (m SecurityPromptModel) Init() tea.Cmd {
	return tea.Batch(
		tea.Tick(time.Millisecond*200, func(t time.Time) tea.Msg {
			return models.TimerTickMsg(t)
		}),
	)
}

func (m SecurityPromptModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case models.TimerTickMsg:
		m.pulseFrame = (m.pulseFrame + 1) % len(pulseColors)
		if m.showAnimation {
			return m, tea.Tick(time.Millisecond*200, func(t time.Time) tea.Msg {
				return models.TimerTickMsg(t)
			})
		}

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, key.NewBinding(key.WithKeys("q", "ctrl+c"))):
			return m, tea.Quit

		case key.Matches(msg, key.NewBinding(key.WithKeys("up", "k"))):
			if m.selectedChoice > 0 {
				m.selectedChoice--
			}

		case key.Matches(msg, key.NewBinding(key.WithKeys("down", "j"))):
			if m.selectedChoice < 1 {
				m.selectedChoice++
			}

		case key.Matches(msg, key.NewBinding(key.WithKeys("1"))):
			m.selectedChoice = 0
			return m, m.handleChoice(true)

		case key.Matches(msg, key.NewBinding(key.WithKeys("2"))):
			m.selectedChoice = 1
			return m, m.handleChoice(false)

		case key.Matches(msg, key.NewBinding(key.WithKeys("enter"))):
			return m, m.handleChoice(m.selectedChoice == 0)
		}
	}

	return m, nil
}

func (m SecurityPromptModel) handleChoice(proceed bool) tea.Cmd {
	return func() tea.Msg {
		if proceed {
			return models.ScreenChangeMsg{
				NewScreen: models.Welcome,
				Data:      m.directory,
			}
		}
		return tea.Quit()
	}
}

func (m SecurityPromptModel) View() string {
	var b strings.Builder

	// Animated title with pulsing border
	animatedTitleStyle := titleStyle.Copy().
		BorderForeground(lipgloss.Color(pulseColors[m.pulseFrame]))

	b.WriteString(animatedTitleStyle.Render("🔒 SECURITY TRUST PROMPT"))
	b.WriteString("\n\n")

	// Directory info with nice styling
	b.WriteString("Aura wants to access the following directory:\n")
	b.WriteString(directoryStyle.Render(fmt.Sprintf("📁 %s", m.directory)))
	b.WriteString("\n\n")

	// Warning box
	warningText := "⚠️  WARNING\n\nThis will allow Aura to:\n• Read and write files in this directory\n• Execute shell commands\n• Access Git repository information\n\nOnly proceed if you trust this codebase!"
	b.WriteString(warningStyle.Render(warningText))
	b.WriteString("\n\n")

	// Choices with selection highlighting
	choices := []string{"✅ Yes, I trust this directory", "❌ No, exit Aura"}
	for i, choice := range choices {
		if i == m.selectedChoice {
			b.WriteString("▶ ")
			b.WriteString(selectedChoiceStyle.Render(fmt.Sprintf("%d. %s", i+1, choice)))
		} else {
			b.WriteString("  ")
			b.WriteString(unselectedChoiceStyle.Render(fmt.Sprintf("%d. %s", i+1, choice)))
		}
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#888888")).Render("Use ↑↓ or 1-2 to select, Enter to confirm, Ctrl+C to quit"))

	return b.String()
}
