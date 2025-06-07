package models

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type Screen int

const (
	SecurityPrompt Screen = iota
	Welcome
	Main
	Configuration
	FileEdit
	CommandExecution
	GitOperation
)

type AppState struct {
	Screen       Screen
	CurrentDir   string
	Config       interface{}
	TrustedDir   bool
	LastError    string
	ProjectFiles []string
	GitStatus    string
}

type ScreenChangeMsg struct {
	NewScreen Screen
	Data      interface{}
}

type ErrorMsg struct {
	Error error
}

type StatusMsg struct {
	Message string
}

type ConfirmationMsg struct {
	Prompt   string
	Command  string
	Callback func(bool) tea.Cmd
}

type FileChangesMsg struct {
	FilePath string
	Changes  string
	Diff     string
}

type GitStatusMsg struct {
	Status string
	Branch string
}

type TimerTickMsg time.Time

type CommandType int

const (
	ReadFile CommandType = iota
	WriteFile
	ExecuteCommand
	GitCommand
)

type PendingCommand struct {
	Type        CommandType
	Description string
	Command     string
	FilePath    string
	Content     string
	Confirmed   bool
}
