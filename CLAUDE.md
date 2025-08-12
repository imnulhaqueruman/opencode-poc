# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

TermAI is a complete AI assistant application built with Go. It features a terminal user interface (TUI) with database persistence, session management, and LLM integration. This represents the evolution from a basic TUI demo to a full AI assistant architecture.

**Key Technologies:**
- Go 1.23.5+
- Bubble Tea (TUI framework)  
- SQLite with migrations (golang-migrate)
- SQLC for type-safe database queries
- Cobra CLI framework
- Viper configuration management

## Development Commands

### Building and Running
```bash
# Run directly
go run .

# Build binary
go build -o termai .

# Build to specific location
go build -o bin/termai .

# Quick run script (includes helpful output)
./run.sh
```

### Development Tasks
```bash
# Install dependencies
go mod download

# Clean dependencies
go mod tidy

# Format code
go fmt ./...

# Vet code for issues
go vet ./...

# Run tests (if any exist)
go test ./...
```

### Database Operations
```bash
# View database schema
sqlite3 ./data/termai.db ".schema"

# List tables
sqlite3 ./data/termai.db ".tables"

# Query sessions
sqlite3 ./data/termai.db "SELECT * FROM sessions;"

# Query messages
sqlite3 ./data/termai.db "SELECT * FROM messages;"
```

### SQLC Code Generation
```bash
# Generate database code (after modifying SQL files)
sqlc generate
```

## Architecture Overview

### Core Services Layer
- **App Context** (`internal/app/`): Dependency injection container with all services
- **Session Service** (`internal/session/`): CRUD operations for chat sessions  
- **Message Service** (`internal/message/`): Message persistence and threading
- **LLM Service** (`internal/llm/`): AI model integration and agent management
- **Permission Service** (`internal/permission/`): Tool execution permissions

### Database Layer (`internal/db/`)
- **SQLite Database**: Single-file database with automatic migrations
- **SQLC Generated Code**: Type-safe query methods in `*.sql.go` files
- **Migrations**: Located in `internal/db/migrations/`
- **SQL Queries**: Raw SQL in `internal/db/sql/` directory

### LLM Integration (`internal/llm/`)
- **Agent System**: ReAct-based agents with tool integration
- **Models**: Provider abstraction (OpenAI, Anthropic, Groq)
- **Tools**: File operations, shell execution, search capabilities
- **Persistent Shell**: Maintains shell state across tool calls

### TUI System (`internal/tui/`)
- **Bento Layout**: Multi-pane interface with resizable sections
- **Components**: Reusable UI elements (editor, sessions list, messages)
- **Pages**: Different application views (REPL, logs)
- **Event-Driven**: Pub/sub system for real-time updates

### Available Tools
The LLM agents have access to these tools:
- `bash`: Execute shell commands with persistent session
- `view`: Read file contents  
- `write`: Create new files
- `edit`: Modify existing files
- `ls`: List directory contents
- `glob`: Pattern-based file searching
- `agent`: Delegate complex tasks to sub-agents

### Configuration System
- **Config File**: `.termai.yaml` in home directory or project root
- **Environment Variables**: Prefixed with `TERMAI_`
- **Defaults**: Set in `cmd/root.go` loadConfig function
- **API Keys**: Configured via viper defaults (should be moved to env vars)

### Pub/Sub Event System (`internal/pubsub/`)
Real-time communication between services:
- Logger events for status updates
- Session events for UI updates  
- Message events for conversation updates
- LLM events for AI responses
- Permission events for user confirmation

## Key File Locations

### Configuration
- `sqlc.yaml` - Database code generation config
- `go.mod` - Go module dependencies

### Database Schema
- `internal/db/migrations/000001_initial.up.sql` - Database schema
- `internal/db/sql/` - Raw SQL queries for SQLC

### Entry Points  
- `main.go` - Application entry point
- `cmd/root.go` - CLI command setup and service initialization

### Core Components
- `internal/app/services.go` - Service dependency injection
- `internal/llm/llm.go` - LLM service implementation
- `internal/tui/tui.go` - Main TUI application

## Important Notes

- **No Tests**: This codebase currently has no test files
- **API Keys Warning**: Current config includes hardcoded API keys that should be moved to environment variables
- **Database**: SQLite database is created automatically in `./data/termai.db`
- **Shell Persistence**: The bash tool maintains shell state between calls
- **Permission System**: Some tools require user permission for execution
- **Memory File**: Look for `termai.md` in working directory for project-specific information

## CLI Usage

```bash
# Start TUI interface
./termai

# Enable debug mode  
./termai --debug

# Show help
./termai --help

# Specify config file
./termai --config=/path/to/config.yaml
```

### TUI Keyboard Shortcuts
- **`?`** - Toggle help overlay
- **`Ctrl+Enter`** - Send message in editor  
- **`Esc`** - Close dialogs/go back
- **`L`** - Switch to logs page
- **`Ctrl+C`** / **`q`** - Quit application