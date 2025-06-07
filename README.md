# Aura - AI-Powered Coding Agent

Aura is a terminal-first, AI-powered coding agent designed to augment the developer workflow. It operates as an interactive Terminal User Interface (TUI) that integrates conversational AI directly into the developer's primary work environment.

## Features

- **Security-First Design**: Trust prompts for new directories with explicit user consent
- **Interactive Terminal Interface**: Beautiful TUI built with Bubble Tea
- **Slash Commands**: Quick access to common operations
- **File System Integration**: Read, analyze, and modify files with AI assistance
- **Command Execution**: Secure shell command execution with confirmation prompts
- **Git Integration**: Built-in Git operations and workflow support
- **Configuration Management**: User-specific settings stored in `~/.aura.conf`

## Installation

### Build from Source

1. Clone or download the project
2. Make sure you have Go 1.19+ installed
3. Build the application:

```bash
go build -o aura .
```

4. Move the binary to your PATH (optional):

```bash
# On Windows
move aura.exe C:\Windows\System32\

# On macOS/Linux
sudo mv aura /usr/local/bin/
```

## Usage

### Starting Aura

Launch Aura in your current directory:

```bash
aura
```

Or specify a directory:

```bash
aura /path/to/your/project
```

### First Run

On first run in a new directory, Aura will show a security prompt:

```
🔒 SECURITY TRUST PROMPT

Aura wants to access: /path/to/your/project

⚠️  WARNING: This will allow Aura to read/write files and execute commands.

1. Yes, proceed
2. No, exit

Choose an option (1-2):
```

Choose option 1 to trust the directory and proceed.

### Main Interface

After the welcome screen, you'll see the main interface:

```
Aura - AI Coding Agent | /path/to/your/project
────────────────────────────────────────────────────────────────────────────────

> Type your request or use /help for commands...
```

### Slash Commands

- `/help` - Show available commands
- `/status` - Display current status and configuration
- `/edit <file>` - Open file for AI-assisted editing
- `/run <command>` - Execute a shell command (with confirmation)
- `/test` - Run project tests
- `/commit` - Interactive Git commit process

### Natural Language Commands

You can also type natural language requests:

- "Create a new file called user.go with a User struct"
- "Add error handling to the main function"
- "Show me the contents of package.json"
- "Run the tests for this project"

### Security Features

- **Directory Trust**: Explicit approval required for new directories
- **Command Confirmation**: All shell commands require user confirmation
- **File Change Confirmation**: File modifications show diffs and require approval

## Configuration

Aura stores configuration in `~/.aura.conf` (YAML format):

```yaml
llm_provider: "anthropic"          # or "openai"
api_key: "your-api-key-here"
trusted_dirs:
  - "/path/to/trusted/project1"
  - "/path/to/trusted/project2"
theme: "default"
default_model: "claude-3-sonnet-20240229"
auto_update: true
```

### Setting up API Keys

Before using AI features, you need to configure your API key:

1. Create or edit `~/.aura.conf`
2. Add your API key:

```yaml
llm_provider: "anthropic"  # or "openai"
api_key: "your-api-key-here"
```

## Requirements

- Go 1.19 or later (for building from source)
- Terminal with color support
- API key for supported LLM provider (Anthropic Claude or OpenAI)

## Supported Platforms

- Windows (with WSL recommended)
- macOS (x86/ARM)
- Linux

## Development

### Project Structure

```
aura/
├── app/           # Application logic and routing
├── cmd/           # Command-line interface
├── config/        # Configuration management
├── models/        # Data models and types
├── ui/            # User interface components
├── main.go        # Entry point
├── go.mod         # Go module definition
└── README.md      # This file
```

### Building

```bash
# Download dependencies
go mod tidy

# Build for current platform
go build -o aura .

# Build for specific platforms
GOOS=linux GOARCH=amd64 go build -o aura-linux-amd64 .
GOOS=darwin GOARCH=amd64 go build -o aura-darwin-amd64 .
GOOS=windows GOARCH=amd64 go build -o aura-windows-amd64.exe .
```

## 🚀 CI/CD Pipeline

Aura uses GitHub Actions for automated testing, code quality checks, and releases:

### Automated Checks
- ✅ **Cross-platform testing** (Linux, Windows, macOS)
- ✅ **Multi-version Go support** (1.20, 1.21)
- ✅ **Comprehensive test suite** (Unit, Integration, E2E, Benchmarks)
- ✅ **Code quality analysis** (golangci-lint)
- ✅ **Security scanning** (gosec)
- ✅ **Code formatting** verification
- ✅ **Dependency management** checks

### Branch Protection
- **`main`** and **`develop`** branches are protected
- All PRs require passing CI checks
- At least one approving review required
- No direct pushes to protected branches

### Automated Releases
- Automatic releases on `main` branch pushes
- Cross-platform binary builds
- Semantic versioning
- Release notes generation

See [CONTRIBUTING.md](CONTRIBUTING.md) for detailed contribution guidelines.

## 🧪 Testing

Aura includes a comprehensive test suite with 145+ tests:

### Running Tests

```bash
# Run all tests with our custom test runner
go run test_runner.go

# Run specific test categories
go test ./config    # Unit tests
go test ./app       # Integration tests
go test ./e2e       # End-to-end tests

# Run with coverage
go test -cover ./...

# Run benchmarks
go test -bench=. ./...
```

### Test Categories

- **Unit Tests (101)**: Test individual components (config, llm, models)
- **Integration Tests (8)**: Test component interactions (app)
- **End-to-End Tests (21)**: Test complete workflows (e2e)
- **Benchmark Tests (15)**: Performance testing

See [TESTING.md](TESTING.md) for detailed testing documentation.

## Contributing

1. Fork the repository
2. Create a feature branch from `develop`
3. Make your changes following our guidelines
4. Add tests for new functionality
5. Ensure all CI checks pass
6. Submit a pull request to `develop`

## License

This project is licensed under the MIT License.

## Security

Aura is designed with security as a primary concern:

- No automatic file modifications without user consent
- No command execution without explicit confirmation
- Directory-based trust model
- All operations are logged and transparent

If you discover a security vulnerability, please report it responsibly.
