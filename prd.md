Of course. Here is a Product Requirements Document (PRD) for "Aura," an AI-powered coding agent in the terminal, based on the provided screenshots and requirements.

---

## Product Requirements Document: Aura - The AI Coding Agent

**Version:** 1.0
**Date:** July 07, 2025
**Author:** SinkaCom AG
**Status:** Draft

### 1. Introduction & Vision

**Aura** is a terminal-first, AI-powered coding agent designed to augment the developer workflow. It operates as an interactive Terminal User Interface (TUI) that integrates conversational AI directly into the developer's primary work environment.

**Vision:** To create the fastest, most intuitive interface for AI-assisted development, allowing developers to build, test, and refactor software without leaving the command line. Aura aims to be a seamless partner in the coding process, understanding context, automating repetitive tasks, and supercharging developer productivity.

### 2. Target Audience

Aura is for **terminal-centric developers** who:
*   Live in the command line and prefer keyboard-driven workflows.
*   Use tools like Vim/Neovim, tmux, Docker, and Git CLI daily.
*   Value speed, efficiency, and minimal context switching.
*   Are comfortable with modern development practices and are open to leveraging AI to improve their workflow.
*   Work on a variety of tasks, from backend services and infrastructure management to frontend development.

### 3. User Problems & Goals

| User Problem | User Goal |
| :--- | :--- |
| Constant context switching between the editor, terminal, and a browser/chat UI for AI assistance disrupts flow and slows down development. | I want to perform all my development tasks, including getting AI help, within a single, unified terminal interface. |
| Writing boilerplate code, unit tests, and documentation is time-consuming and repetitive. | I want to automate the generation of repetitive code and documentation with simple, natural language commands. |
| It's difficult to quickly understand a new or unfamiliar codebase. | I want to get high-level summaries of files, directories, or the entire project to accelerate onboarding and feature development. |
| Complex Git operations or shell commands require looking up documentation. | I want to execute complex commands using natural language, without having to remember the exact syntax. |
| Refactoring code or applying changes across multiple files is tedious and error-prone. | I want an intelligent agent that can understand my refactoring goals and apply the necessary changes safely across the project. |

### 4. Core Features

#### F1: Onboarding & Project Scoping

*   **Invocation:** Aura is launched by running the `aura` command in a project directory. If run without a directory argument, it will use the current working directory.
*   **Security Trust Prompt (Screen 1):**
    *   On the first run in a new, untrusted directory, Aura will present a security prompt.
    *   It will clearly state the directory path it needs access to.
    *   It will warn the user about the risks of reading and executing files from an untrusted source.
    *   It will provide a clear choice: `1. Yes, proceed` or `2. No, exit`.
    *   Trusted directories will be stored in the configuration file to avoid repeated prompts.
*   **Welcome Screen (Screen 2):**
    *   Upon successful launch, Aura displays a welcome message.
    *   It shows the current working directory (`cwd`).
    *   It provides helpful "getting started" tips and key commands like `/help` and `/status`.

#### F2: Conversational AI Core

*   **Interactive Prompt:** The core of the UI is a persistent input prompt (`>`) where the user can type natural language requests or slash commands.
*   **Context Awareness:** Aura maintains context of the current project, including the file tree and the content of recently viewed/edited files.
*   **Slash Commands:** Structured commands for common operations, providing a more deterministic way to interact with the agent.
    *   `/help`: Displays a list of available commands and tips.
    *   `/status`: Shows the current project, LLM model in use, and other relevant status information.
    *   `/edit <filepath>`: Opens a file for AI-driven editing.
    *   `/run <command>`: Executes a shell command (with confirmation).
    *   `/test`: Runs the project's test suite.
    *   `/commit`: Initiates an interactive Git commit process.

#### F3: File System & Code Interaction

*   **File Analysis:** Users can ask Aura to read and analyze files or directories (e.g., `"Summarize src/api/auth.go"`, `"What are the main dependencies in package.json?"`).
*   **Code Generation & Editing:** Aura can create new files or modify existing ones based on user prompts.
    *   Example: `"Create a new file named user_service.go with a UserService struct and CRUD functions."`
    *   Example (after `/edit user_service.go`): `"Add error handling to the UpdateUser function."`
*   **Diff & Confirmation:** Before writing any changes to disk, Aura must show the user a diff of the proposed changes and require explicit confirmation.

#### F4: Secure Command Execution

*   **Natural Language to Shell:** Users can ask Aura to perform actions that require shell commands (e.g., `"Install the dependencies"`, `"Run the linter"`).
*   **Execution Confirmation:** For any command that Aura intends to execute, it **must** first display the exact command to the user and ask for confirmation (`Aura wants to run: 'npm install'. Proceed? [y/N]`). This is a critical security feature.

#### F5: Version Control (Git) Integration

*   **Git Awareness:** Aura can read the Git status, view diffs, and list branches.
*   **Automated Git Operations:** Users can ask Aura to perform Git operations.
    *   `"Create a new branch called feature/user-auth"`
    *   `"Stage all changes in the src/ directory"`
    *   `"Commit the staged changes with the message 'feat: Implement user authentication endpoint'"`

#### F6: System & Configuration

*   **Configuration File:** All user-specific settings are stored in `~/.aura.conf`. This includes:
    *   LLM provider and API key (e.g., Anthropic, OpenAI).
    *   List of trusted project directories.
    *   Theme/color preferences.
    *   Default model selection.
*   **Auto-Update Mechanism:** Aura will have a built-in mechanism to check for new versions and prompt the user to update. It will handle update failures gracefully, as suggested by the error message in Screen 2.

### 5. User Flow Example

1.  **Developer `cd`s into their project:** `cd ~/projects/my-app`
2.  **Launches Aura:** `$ aura`
3.  **First-time use:** Aura displays the security prompt for `~/projects/my-app`. The developer types `1` and presses Enter to trust it.
4.  **Main Interface:** The welcome screen appears, showing `cwd: ~/projects/my-app`.
5.  **Task:** The developer wants to add a new API endpoint.
6.  **Developer types:** `> create a new git branch named feature/add-health-check`
7.  **Aura responds:** `Switched to a new branch 'feature/add-health-check'`.
8.  **Developer types:** `> create a new file router/health.go and add a basic http handler for a GET /health endpoint that returns a 200 OK with a JSON body {"status": "ok"}`
9.  **Aura responds:**
    ```
    I will create the file router/health.go with the following content:
    <shows generated Go code>
    
    Proceed? [Y/n]
    ```
10. **Developer confirms:** `y`
11. **Developer types:** `> now, register this new handler in main.go`
12. **Aura responds:**
    ```
    I will apply the following changes to main.go:
    
    --- a/main.go
    +++ b/main.go
    <shows diff of changes>
    
    Apply changes? [Y/n]
    ```
13. **Developer confirms:** `y`
14. **Developer types:** `> /test`
15. **Aura confirms and runs tests:** `Aura wants to run: 'go test ./...'. Proceed? [y/N]` -> `y` -> (shows test output).
16. **Developer types:** `> /commit`
17. **Aura stages all changes and asks for a commit message:** `AI-generated commit message: 'feat: Add /health endpoint'. Use this message? [Y/n]`
18. **Developer exits:** `Ctrl+C`

### 6. Non-Functional Requirements

| Requirement | Description |
| :--- | :--- |
| **Performance** | The TUI must be highly responsive with minimal input latency. AI responses should be streamed to the user to avoid long waits. |
| **Security** | This is paramount. The trust model, explicit confirmation for file writes, and mandatory confirmation for command execution are non-negotiable. The agent must not execute code without user consent. |
| **Platform Support** | The application must be a single, statically-linked binary compiled for Linux, macOS (x86/ARM), and Windows (via WSL). |
| **Reliability** | The application should handle API errors, network issues, and unexpected command outputs gracefully without crashing. |
| **Usability** | The interface should be intuitive for developers familiar with the command line. Keyboard shortcuts for common actions should be available and discoverable via `/help`. |

### 7. Technical Implementation Details

*   **Programming Language:** **Go**
*   **TUI Framework:** A robust Go TUI library like **Bubble Tea** is recommended for its component-based model, which is well-suited for this type of application.
*   **Configuration:** The `~/.aura.conf` file will be parsed using a standard library or a simple TOML/YAML parser.
*   **Backend:** Aura is a frontend for a large language model (LLM). It will interact with LLM provider APIs (e.g., Anthropic, OpenAI) via REST. The specific provider will be user-configurable.

### 8. Success Metrics

*   **Adoption:** Number of unique daily and weekly active users.
*   **Engagement:** Average session duration; number of commands/interactions per session.
*   **Task Success Rate:** Percentage of user requests that result in a successful file modification, command execution, or commit.
*   **Qualitative Feedback:** User feedback collected via GitHub issues, surveys, and community channels.
*   **Retention:** Week 1 and Month 1 user retention rate.

### 9. Future Scope (v2.0 and beyond)

*   **Deeper IDE Features:** Integrate with Language Server Protocol (LSP) for code completion, diagnostics, and go-to-definition within the TUI.
*   **Interactive Debugging:** Allow users to step through code with a debugger controlled by natural language.
*   **Plugin Architecture:** Allow third-party developers to extend Aura's functionality with new commands and integrations.
*   **Multi-Agent Mode:** Allow the user to launch specialized agents (e.g., a "Test Agent", a "DBA Agent") that can collaborate on a task.