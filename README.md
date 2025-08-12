# TermAI - Evolution to Full AI Assistant (OpenCode Genesis)

This is the **evolution** of TermAI from a basic TUI demo to a full AI assistant with database persistence - tracing the journey that would eventually become OpenCode.

## Current State Analysis (After "add initial stuff" commit)

**Project Name**: `termai` (Terminal AI)  
**Language**: Go  
**Architecture**: Full-stack terminal AI assistant with SQLite database  
**Author**: Kujtim Hoxha  
**Integration**: Initial commit → Help system → Full AI assistant architecture

## Architecture Overview

This evolved version is a **complete AI assistant application** with database persistence and session management:

### Major Architectural Changes

**Database Layer**:
- SQLite database with migrations
- Session persistence with message history
- SQLC code generation for type-safe queries

**CLI Framework**:
- Cobra CLI with subcommands and configuration
- Viper for configuration management
- Professional command-line interface

**Session Management**:
- Create, read, update, delete sessions
- Message threading with user/assistant roles
- Real-time session updates via pub/sub

**Enhanced TUI**:
- Multi-pane layout with sessions list
- Message display with threading
- Enhanced editor with message sending
- Status notifications and error handling

### Core Components

1. **CLI Application** (`cmd/root.go`)
   - Cobra CLI framework with configuration
   - Database initialization and migrations
   - Application context with services

2. **Database Layer** (`internal/db/`)
   - SQLite with migrations using golang-migrate
   - SQLC generated code for type-safe queries
   - Session and message persistence

3. **Session Service** (`internal/session/`)
   - CRUD operations for sessions and messages
   - Pub/sub events for real-time updates
   - JSON message serialization

4. **TUI System** (`internal/tui/`)
   - **Pages**: Enhanced REPL with session management
   - **Components**: Sessions list, Messages display, Enhanced editor
   - **Layouts**: Bento layout with help/dialog overlays

5. **App Context** (`internal/app/`)
   - Service dependency injection
   - Database connection management
   - Configuration handling

### Key Features (Current)

- **Session Management**: Create, select, and manage AI conversation sessions
- **Message Threading**: Full conversation history with user/assistant messages
- **Database Persistence**: SQLite-backed data storage with migrations
- **Real-time Updates**: Pub/sub system for live session updates
- **Enhanced UI**: Multi-pane layout with sessions list and message display
- **Professional CLI**: Cobra-based command system with configuration

## Dependencies (Go Modules)

### Core Framework
- **Cobra**: Professional CLI application framework
- **Viper**: Configuration management
- **Bubble Tea**: TUI framework for terminal applications
- **Bubbles**: Pre-built UI components (list, viewport, textarea)
- **Lipgloss**: Terminal styling and layout
- **Huh**: Enhanced forms and dialogs

### Database & Persistence  
- **SQLite3**: Embedded database (github.com/mattn/go-sqlite3)
- **golang-migrate**: Database migrations
- **UUID**: Unique ID generation for sessions/messages

### Additional
- **Glamour**: Markdown rendering for terminals
- **Catppuccin**: Color theme support

## Setup Instructions

### Prerequisites

- **Go 1.23.5+** (specified in go.mod)
- Terminal with true color support (recommended)

### Installation & Running

1. **Navigate to the project directory**:
   ```bash
   cd /home/poridhi/development/opencode-poc
   ```

2. **Install dependencies**:
   ```bash
   go mod download
   ```

3. **Initialize the database**:
   ```bash
   # The application will automatically create and migrate the database on first run
   ```

4. **Build and run**:
   ```bash
   # Build the application
   go build -o termai .
   
   # Run the application  
   ./termai
   
   # Or run directly
   go run .
   ```

5. **CLI Usage**:
   ```bash
   # Show help
   ./termai --help
   
   # Run in debug mode
   ./termai --debug
   ```

### Development Commands

```bash
# Run in development mode
go run .

# Build binary
go build -o bin/termai .

# Clean build
go mod tidy && go build .

# Run tests
go test ./...

# Format code
go fmt ./...

# Check for issues
go vet ./...

# View database (SQLite CLI)
sqlite3 ./data/termai.db ".tables"
sqlite3 ./data/termai.db "SELECT * FROM sessions;"
```

## Interface Overview

The current TUI features a **professional AI assistant interface**:

### Main Interface (REPL Page)
- **Left Pane**: Sessions list with create/select functionality
- **Right Top**: Message display with conversation threading  
- **Right Bottom**: Enhanced editor with send functionality (Ctrl+Enter)

### Session Management
- **New Session**: Create sessions with custom titles
- **Session Selection**: Click or navigate to switch sessions
- **Message Threading**: Full conversation history per session
- **Real-time Updates**: Live session updates via pub/sub

### Database Integration
- **Persistent Sessions**: All sessions stored in SQLite
- **Message History**: Complete conversation threading
- **Automatic Migrations**: Database schema managed automatically

## Key Architecture Decisions

### 1. **Bubble Tea Framework**
- Event-driven TUI architecture
- Component-based design
- Message passing between components

### 2. **Bento Layout System**
- Multi-pane interface design
- Resizable and responsive layouts
- Modular pane management

### 3. **Pub/Sub Logging**
- Real-time log streaming to UI
- Decoupled logging from display
- Subscription-based event handling

### 4. **Page-based Navigation**
- Simple page routing system
- State preservation between pages
- Lazy page initialization

## Evolution Notes

This version (after "add help" commit) includes:

### ✅ Added in Help Commit:
- **Help System** (`?` key) - Toggle contextual help overlay
- **Status Bar** - Shows version, help hint, errors/info messages
- **Version Management** - Build-time version tracking
- **Enhanced Keybindings** - More navigation options
- **Error/Info Messaging** - User feedback system
- **Improved Layout** - Status bar integration

### Still Missing (Future Commits):
- **No AI integration yet** - Pure TUI framework
- **No file operations** - Just UI structure  
- **No provider system** - Missing AI/LLM integration
- **No tools** - No bash, edit, read, etc. tools
- **Go-only** - Before the TypeScript/Bun rewrite

The architecture shows early design decisions that would influence the final product:
- Multi-pane interfaces with help overlay
- Event-driven architecture with status messaging
- Component separation and modularity
- Real-time logging/messaging with user feedback

## Next Steps in Evolution

To see how this evolved into OpenCode, you would trace through commits to see:

1. **AI Integration** - Adding LLM providers
2. **Tool System** - File operations, shell commands
3. **Language Migration** - Go → TypeScript/Bun
4. **HTTP API** - Server mode development
5. **Multi-provider Support** - Anthropic, OpenAI, etc.

## Running the Current Version

This version includes a **complete AI assistant foundation** with:
- Full session management with database persistence
- Multi-pane layout with enhanced UI components
- Professional CLI framework with configuration
- Real-time message threading and updates
- Status notifications and error handling
- Database migrations and data management

### Keyboard Shortcuts:
- **`?`** - Toggle help overlay (shows all available commands)
- **`Ctrl+Enter`** - Send message in editor
- **`Esc`** - Close dialogs or go back
- **`L`** - Switch to logs page  
- **`Ctrl+C` / `q`** - Quit application (with confirmation dialog)

### CLI Commands:
```bash
# Basic usage
./termai              # Start TUI interface
./termai --debug      # Enable debug mode
./termai --help       # Show help

# Configuration
./termai --config=/path/to/config.yaml
```

### Database Structure:
```sql
-- Sessions table
CREATE TABLE sessions (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    messages TEXT NOT NULL DEFAULT '[]',  -- JSON array
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

This version represents the **complete architecture foundation** for an AI assistant - with persistent sessions, message threading, and professional CLI framework. The next steps would be adding actual AI provider integration (OpenAI, Anthropic, etc.) to complete the assistant functionality.
```html
<iframe width="560" height="315" src="[https://www.youtube.com/embed/abc123"](https://drive.google.com/file/d/1eMCPuRoFbhSLIld27N0bNkq_hKD1HEA9/view?usp=sharing" frameborder="0" allowfullscreen></iframe>
```
