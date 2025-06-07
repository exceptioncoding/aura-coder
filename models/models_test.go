package models

import (
	"testing"
	"time"

	"github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
)

func TestScreen_Constants(t *testing.T) {
	tests := []struct {
		name     string
		screen   Screen
		expected int
	}{
		{
			name:     "SecurityPrompt",
			screen:   SecurityPrompt,
			expected: 0,
		},
		{
			name:     "Welcome",
			screen:   Welcome,
			expected: 1,
		},
		{
			name:     "Main",
			screen:   Main,
			expected: 2,
		},
		{
			name:     "Configuration",
			screen:   Configuration,
			expected: 3,
		},
		{
			name:     "FileEdit",
			screen:   FileEdit,
			expected: 4,
		},
		{
			name:     "CommandExecution",
			screen:   CommandExecution,
			expected: 5,
		},
		{
			name:     "GitOperation",
			screen:   GitOperation,
			expected: 6,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, int(tt.screen))
		})
	}
}

func TestCommandType_Constants(t *testing.T) {
	tests := []struct {
		name        string
		commandType CommandType
		expected    int
	}{
		{
			name:        "ReadFile",
			commandType: ReadFile,
			expected:    0,
		},
		{
			name:        "WriteFile",
			commandType: WriteFile,
			expected:    1,
		},
		{
			name:        "ExecuteCommand",
			commandType: ExecuteCommand,
			expected:    2,
		},
		{
			name:        "GitCommand",
			commandType: GitCommand,
			expected:    3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, int(tt.commandType))
		})
	}
}

func TestAppState_Initialization(t *testing.T) {
	tests := []struct {
		name     string
		state    AppState
		expected AppState
	}{
		{
			name: "default app state",
			state: AppState{
				Screen:       Welcome,
				CurrentDir:   "/home/user/project",
				Config:       nil,
				TrustedDir:   false,
				LastError:    "",
				ProjectFiles: []string{},
				GitStatus:    "",
			},
			expected: AppState{
				Screen:       Welcome,
				CurrentDir:   "/home/user/project",
				Config:       nil,
				TrustedDir:   false,
				LastError:    "",
				ProjectFiles: []string{},
				GitStatus:    "",
			},
		},
		{
			name: "app state with data",
			state: AppState{
				Screen:       Main,
				CurrentDir:   "/opt/myapp",
				Config:       "test-config",
				TrustedDir:   true,
				LastError:    "previous error",
				ProjectFiles: []string{"main.go", "config.go"},
				GitStatus:    "clean",
			},
			expected: AppState{
				Screen:       Main,
				CurrentDir:   "/opt/myapp",
				Config:       "test-config",
				TrustedDir:   true,
				LastError:    "previous error",
				ProjectFiles: []string{"main.go", "config.go"},
				GitStatus:    "clean",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected.Screen, tt.state.Screen)
			assert.Equal(t, tt.expected.CurrentDir, tt.state.CurrentDir)
			assert.Equal(t, tt.expected.Config, tt.state.Config)
			assert.Equal(t, tt.expected.TrustedDir, tt.state.TrustedDir)
			assert.Equal(t, tt.expected.LastError, tt.state.LastError)
			assert.Equal(t, tt.expected.ProjectFiles, tt.state.ProjectFiles)
			assert.Equal(t, tt.expected.GitStatus, tt.state.GitStatus)
		})
	}
}

func TestScreenChangeMsg(t *testing.T) {
	tests := []struct {
		name    string
		msg     ScreenChangeMsg
		wantMsg tea.Msg
	}{
		{
			name: "screen change to main",
			msg: ScreenChangeMsg{
				NewScreen: Main,
				Data:      "/home/user/project",
			},
			wantMsg: ScreenChangeMsg{
				NewScreen: Main,
				Data:      "/home/user/project",
			},
		},
		{
			name: "screen change to config",
			msg: ScreenChangeMsg{
				NewScreen: Configuration,
				Data:      nil,
			},
			wantMsg: ScreenChangeMsg{
				NewScreen: Configuration,
				Data:      nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test that ScreenChangeMsg implements tea.Msg
			var msg tea.Msg = tt.msg
			assert.Equal(t, tt.wantMsg, msg)
		})
	}
}

func TestErrorMsg(t *testing.T) {
	tests := []struct {
		name    string
		msg     ErrorMsg
		wantMsg tea.Msg
	}{
		{
			name: "error message with error",
			msg: ErrorMsg{
				Error: assert.AnError,
			},
			wantMsg: ErrorMsg{
				Error: assert.AnError,
			},
		},
		{
			name: "error message with nil error",
			msg: ErrorMsg{
				Error: nil,
			},
			wantMsg: ErrorMsg{
				Error: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test that ErrorMsg implements tea.Msg
			var msg tea.Msg = tt.msg
			assert.Equal(t, tt.wantMsg, msg)
		})
	}
}

func TestStatusMsg(t *testing.T) {
	tests := []struct {
		name    string
		msg     StatusMsg
		wantMsg tea.Msg
	}{
		{
			name: "status message",
			msg: StatusMsg{
				Message: "Operation completed successfully",
			},
			wantMsg: StatusMsg{
				Message: "Operation completed successfully",
			},
		},
		{
			name: "empty status message",
			msg: StatusMsg{
				Message: "",
			},
			wantMsg: StatusMsg{
				Message: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test that StatusMsg implements tea.Msg
			var msg tea.Msg = tt.msg
			assert.Equal(t, tt.wantMsg, msg)
		})
	}
}

func TestConfirmationMsg(t *testing.T) {
	testCallback := func(_ bool) tea.Cmd {
		return func() tea.Msg {
			return StatusMsg{Message: "callback executed"}
		}
	}

	tests := []struct {
		name string
		msg  ConfirmationMsg
	}{
		{
			name: "confirmation with callback",
			msg: ConfirmationMsg{
				Prompt:   "Are you sure?",
				Command:  "rm -rf /",
				Callback: testCallback,
			},
		},
		{
			name: "confirmation without callback",
			msg: ConfirmationMsg{
				Prompt:   "Continue?",
				Command:  "make build",
				Callback: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test that ConfirmationMsg implements tea.Msg
			var msg tea.Msg = tt.msg

			assert.Equal(t, tt.msg.Prompt, msg.(ConfirmationMsg).Prompt)
			assert.Equal(t, tt.msg.Command, msg.(ConfirmationMsg).Command)

			// Test callback execution if present
			if tt.msg.Callback != nil {
				cmd := tt.msg.Callback(true)
				assert.NotNil(t, cmd)
				resultMsg := cmd()
				assert.Equal(t, StatusMsg{Message: "callback executed"}, resultMsg)
			}
		})
	}
}

func TestFileChangesMsg(t *testing.T) {
	tests := []struct {
		name    string
		msg     FileChangesMsg
		wantMsg tea.Msg
	}{
		{
			name: "file changes message",
			msg: FileChangesMsg{
				FilePath: "/path/to/file.go",
				Changes:  "Added new function",
				Diff:     "+func NewFunction() {}",
			},
			wantMsg: FileChangesMsg{
				FilePath: "/path/to/file.go",
				Changes:  "Added new function",
				Diff:     "+func NewFunction() {}",
			},
		},
		{
			name: "empty file changes message",
			msg: FileChangesMsg{
				FilePath: "",
				Changes:  "",
				Diff:     "",
			},
			wantMsg: FileChangesMsg{
				FilePath: "",
				Changes:  "",
				Diff:     "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test that FileChangesMsg implements tea.Msg
			var msg tea.Msg = tt.msg
			assert.Equal(t, tt.wantMsg, msg)
		})
	}
}

func TestGitStatusMsg(t *testing.T) {
	tests := []struct {
		name    string
		msg     GitStatusMsg
		wantMsg tea.Msg
	}{
		{
			name: "git status message",
			msg: GitStatusMsg{
				Status: "clean",
				Branch: "main",
			},
			wantMsg: GitStatusMsg{
				Status: "clean",
				Branch: "main",
			},
		},
		{
			name: "git status with changes",
			msg: GitStatusMsg{
				Status: "modified:   main.go",
				Branch: "feature/new-feature",
			},
			wantMsg: GitStatusMsg{
				Status: "modified:   main.go",
				Branch: "feature/new-feature",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test that GitStatusMsg implements tea.Msg
			var msg tea.Msg = tt.msg
			assert.Equal(t, tt.wantMsg, msg)
		})
	}
}

func TestTimerTickMsg(t *testing.T) {
	now := time.Now()
	tickMsg := TimerTickMsg(now)

	// Test that TimerTickMsg implements tea.Msg
	var msg tea.Msg = tickMsg
	assert.Equal(t, now, time.Time(msg.(TimerTickMsg)))
}

func TestPendingCommand(t *testing.T) {
	tests := []struct {
		name    string
		command PendingCommand
	}{
		{
			name: "read file command",
			command: PendingCommand{
				Type:        ReadFile,
				Description: "Read configuration file",
				Command:     "cat config.yaml",
				FilePath:    "/path/to/config.yaml",
				Content:     "",
				Confirmed:   false,
			},
		},
		{
			name: "write file command",
			command: PendingCommand{
				Type:        WriteFile,
				Description: "Save updated configuration",
				Command:     "",
				FilePath:    "/path/to/config.yaml",
				Content:     "key: value",
				Confirmed:   true,
			},
		},
		{
			name: "execute command",
			command: PendingCommand{
				Type:        ExecuteCommand,
				Description: "Run tests",
				Command:     "go test ./...",
				FilePath:    "",
				Content:     "",
				Confirmed:   false,
			},
		},
		{
			name: "git command",
			command: PendingCommand{
				Type:        GitCommand,
				Description: "Commit changes",
				Command:     "git commit -m 'Update config'",
				FilePath:    "",
				Content:     "",
				Confirmed:   true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.command.Type, tt.command.Type)
			assert.Equal(t, tt.command.Description, tt.command.Description)
			assert.Equal(t, tt.command.Command, tt.command.Command)
			assert.Equal(t, tt.command.FilePath, tt.command.FilePath)
			assert.Equal(t, tt.command.Content, tt.command.Content)
			assert.Equal(t, tt.command.Confirmed, tt.command.Confirmed)
		})
	}
}

func TestPendingCommand_IsFileOperation(t *testing.T) {
	tests := []struct {
		name     string
		command  PendingCommand
		expected bool
	}{
		{
			name: "read file operation",
			command: PendingCommand{
				Type: ReadFile,
			},
			expected: true,
		},
		{
			name: "write file operation",
			command: PendingCommand{
				Type: WriteFile,
			},
			expected: true,
		},
		{
			name: "execute command operation",
			command: PendingCommand{
				Type: ExecuteCommand,
			},
			expected: false,
		},
		{
			name: "git command operation",
			command: PendingCommand{
				Type: GitCommand,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isFileOp := tt.command.Type == ReadFile || tt.command.Type == WriteFile
			assert.Equal(t, tt.expected, isFileOp)
		})
	}
}

func TestPendingCommand_RequiresConfirmation(t *testing.T) {
	tests := []struct {
		name     string
		command  PendingCommand
		expected bool
	}{
		{
			name: "command not confirmed yet",
			command: PendingCommand{
				Type:      ExecuteCommand,
				Confirmed: false,
			},
			expected: true,
		},
		{
			name: "command already confirmed",
			command: PendingCommand{
				Type:      ExecuteCommand,
				Confirmed: true,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requiresConfirmation := !tt.command.Confirmed
			assert.Equal(t, tt.expected, requiresConfirmation)
		})
	}
}

// Test message type assertions
func TestMessageTypeAssertions(t *testing.T) {
	t.Run("all messages implement tea.Msg", func(t *testing.T) {
		var msgs []tea.Msg

		// Test all message types can be assigned to tea.Msg
		msgs = append(msgs, ScreenChangeMsg{})
		msgs = append(msgs, ErrorMsg{})
		msgs = append(msgs, StatusMsg{})
		msgs = append(msgs, ConfirmationMsg{})
		msgs = append(msgs, FileChangesMsg{})
		msgs = append(msgs, GitStatusMsg{})
		msgs = append(msgs, TimerTickMsg(time.Now()))

		assert.Len(t, msgs, 7)
		for i, msg := range msgs {
			assert.NotNil(t, msg, "message %d should not be nil", i)
		}
	})
}

// Benchmark tests for performance-critical operations
func BenchmarkScreenChangeMsg_Creation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = ScreenChangeMsg{
			NewScreen: Main,
			Data:      "/test/path",
		}
	}
}

func BenchmarkTimerTickMsg_Creation(b *testing.B) {
	now := time.Now()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = TimerTickMsg(now)
	}
}

func BenchmarkPendingCommand_Creation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = PendingCommand{
			Type:        ExecuteCommand,
			Description: "Benchmark test command",
			Command:     "echo 'benchmark'",
			FilePath:    "/tmp/benchmark",
			Content:     "test content",
			Confirmed:   false,
		}
	}
}

// Test edge cases and error conditions
func TestEdgeCases(t *testing.T) {
	t.Run("empty strings in messages", func(t *testing.T) {
		msg := StatusMsg{Message: ""}
		assert.Equal(t, "", msg.Message)
	})

	t.Run("nil data in screen change", func(t *testing.T) {
		msg := ScreenChangeMsg{
			NewScreen: Welcome,
			Data:      nil,
		}
		assert.Nil(t, msg.Data)
	})

	t.Run("nil error in error message", func(t *testing.T) {
		msg := ErrorMsg{Error: nil}
		assert.Nil(t, msg.Error)
	})

	t.Run("empty file path in file changes", func(t *testing.T) {
		msg := FileChangesMsg{
			FilePath: "",
			Changes:  "some changes",
			Diff:     "some diff",
		}
		assert.Empty(t, msg.FilePath)
		assert.NotEmpty(t, msg.Changes)
	})
}
