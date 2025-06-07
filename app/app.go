package app

import (
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"

	"aura/config"
	"aura/models"
	"aura/ui"
)

type App struct {
	program    *tea.Program
	config     *config.Config
	currentDir string
	state      models.AppState
}

func NewApp() (*App, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	currentDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	currentDir, err = filepath.Abs(currentDir)
	if err != nil {
		return nil, err
	}

	app := &App{
		config:     cfg,
		currentDir: currentDir,
		state: models.AppState{
			CurrentDir: currentDir,
			Config:     cfg,
			TrustedDir: cfg.IsTrustedDir(currentDir),
		},
	}

	var initialModel tea.Model
	// Check if configuration is complete (has API key)
	switch {
	case cfg.APIKey == "" || cfg.LLMProvider == "":
		app.state.Screen = models.Configuration
		initialModel = ui.NewConfigModel(currentDir)
	case app.state.TrustedDir:
		app.state.Screen = models.Welcome
		initialModel = ui.NewWelcomeModel(currentDir)
	default:
		app.state.Screen = models.SecurityPrompt
		initialModel = ui.NewSecurityPromptModel(currentDir)
	}

	rootModel := NewRootModel(initialModel, app)
	app.program = tea.NewProgram(rootModel, tea.WithAltScreen())

	return app, nil
}

func (a *App) Run() error {
	_, err := a.program.Run()
	return err
}

func (a *App) GetConfig() *config.Config {
	return a.config
}

func (a *App) GetCurrentDir() string {
	return a.currentDir
}

func (a *App) AddTrustedDir(dir string) error {
	return a.config.AddTrustedDir(dir)
}
