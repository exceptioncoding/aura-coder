package app

import (
	tea "github.com/charmbracelet/bubbletea"

	"aura/models"
	"aura/ui"
)

type RootModel struct {
	currentModel tea.Model
	app          *App
}

func NewRootModel(initialModel tea.Model, app *App) RootModel {
	return RootModel{
		currentModel: initialModel,
		app:          app,
	}
}

func (m RootModel) Init() tea.Cmd {
	return m.currentModel.Init()
}

func (m RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if screenMsg, ok := msg.(models.ScreenChangeMsg); ok {
		return m.handleScreenChange(screenMsg)
	}

	var cmd tea.Cmd
	m.currentModel, cmd = m.currentModel.Update(msg)
	return m, cmd
}

func (m RootModel) handleScreenChange(msg models.ScreenChangeMsg) (RootModel, tea.Cmd) {
	switch msg.NewScreen {
	case models.Welcome:
		if data, ok := msg.Data.(string); ok {
			_ = m.app.AddTrustedDir(data)
			m.currentModel = ui.NewWelcomeModel(data)
		}

	case models.Main:
		if data, ok := msg.Data.(string); ok {
			m.currentModel = ui.NewMainModel(data)
		}

	case models.SecurityPrompt:
		if data, ok := msg.Data.(string); ok {
			m.currentModel = ui.NewSecurityPromptModel(data)
		}

	case models.Configuration:
		if data, ok := msg.Data.(string); ok {
			m.currentModel = ui.NewConfigModel(data)
		}
	}

	return m, m.currentModel.Init()
}

func (m RootModel) View() string {
	return m.currentModel.View()
}
