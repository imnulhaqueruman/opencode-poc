
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

## 🚀 Quick Start Guide

### Prerequisites

- **Go 1.23.5+** (check with `go version`)
- **Git** for cloning the repository
- **Terminal** with true color support (recommended)
- **API Key** for at least one LLM provider (see configuration below)

### Installation & Setup

#### 1. **Clone the Repository**
```bash
git clone https://github.com/imnulhaqueruman/opencode-poc.git
cd opencode-poc
```

#### 2. **Install Dependencies**
```bash
# Download all Go dependencies
go mod download

# Clean up dependencies (if needed)
go mod tidy
```

#### 3. **Configure LLM Provider**
Create a `.termai.yaml` file in the project root with your API configuration:

```yaml
# .termai.yaml - Configuration file
data:
    directory: .termai          # Local data storage directory

log:
    level: info                 # Log level: debug, info, warn, error

model:
    coder: claude-3-5-sonnet-20241022    # Primary coding model
    task: claude-3-5-sonnet-20241022     # Task execution model
    coderMaxTokens: 8000
    taskMaxTokens: 4000

providers:
    # Anthropic Claude (Recommended)
    anthropic:
        apiKey: "your-anthropic-api-key-here"
        enabled: true
    
    # OpenAI GPT
    openai:
        apiKey: "your-openai-api-key-here" 
        enabled: false
    
    # Google Gemini
    gemini:
        apiKey: "your-gemini-api-key-here"
        enabled: false
        
    # Groq (Fast inference)
    groq:
        apiKey: "your-groq-api-key-here"
        enabled: false
```

**Alternative: Environment Variables**
```bash
# Set API keys via environment variables
export ANTHROPIC_API_KEY="your-anthropic-api-key"
export OPENAI_API_KEY="your-openai-api-key"
export GEMINI_API_KEY="your-gemini-api-key"
export GROQ_API_KEY="your-groq-api-key"
```

#### 4. **Build and Run**
```bash
# Option 1: Run directly (for development)
go run .

# Option 2: Build and run executable
go build -o termai .
./termai

```

### 🎯 Getting API Keys

#### **Anthropic Claude (Recommended)**
1. Visit [console.anthropic.com](https://console.anthropic.com)
2. Create account and get API key
3. Models: `claude-3-5-sonnet-20241022`, `claude-3-haiku-20240307`

#### **OpenAI**
1. Visit [platform.openai.com](https://platform.openai.com)
2. Generate API key in API section
3. Models: `gpt-4o`, `gpt-4o-mini`, `gpt-3.5-turbo`

#### **Google Gemini**
1. Visit [ai.google.dev](https://ai.google.dev)
2. Get API key from Google AI Studio
3. Models: `gemini-2.0-flash`, `gemini-1.5-pro`

#### **Groq (Fast & Free)**
1. Visit [console.groq.com](https://console.groq.com)
2. Sign up and get free API key
3. Models: `qwen-qwq`, `llama-3.1-70b-versatile`

### 🎮 Usage Instructions

#### **First Run**
```bash
# Start the application
./termai

# Or with debug logging
./termai --debug
```
[watch](./resources/puku.gif)

#### **Interface Overview**
```
┌─────────────────────────────────────────────────────────────────┐
│                        TermAI Interface                        │
├─────────────────┬───────────────────────────────────────────────┤
│   Sessions      │              Messages                         │
│                 │                                               │
│ • New Session   │  User: How do I create a Go project?         │
│ • Project Help  │                                               │
│ • Debug Issue   │  Assistant: I'll help you create a Go        │
│                 │  project. First, let me check if Go is       │
│                 │  installed...                                 │
│                 │                                               │
├─────────────────┼───────────────────────────────────────────────┤
│                 │              Editor                           │
│                 │                                               │
│                 │  Type your message here...                    │
│                 │  Ctrl+Enter to send                           │
└─────────────────┴───────────────────────────────────────────────┘
```

#### **Keyboard Shortcuts**
- **`?`** - Toggle help overlay
- **`Ctrl+Enter`** - Send message 
- **`Esc`** - Close dialogs/go back
- **`L`** - Switch to logs page
- **`Ctrl+C`** / **`q`** - Quit application

#### **CLI Options**
```bash
./termai --help              # Show all available options
./termai --debug             # Enable debug logging
./termai --config=/path/to/config.yaml  # Custom config file
```

## 🛠️ Development Commands

```bash
# Development workflow
go run .                     # Run directly for development
go build -o termai .         # Build binary
go mod tidy                  # Clean dependencies
go mod download              # Download dependencies

# Code quality
go fmt ./...                 # Format code
go vet ./...                 # Static analysis
go test ./...                # Run tests

# Database operations
sqlite3 ./.termai/termai.db ".tables"                    # List tables
sqlite3 ./.termai/termai.db "SELECT * FROM sessions;"    # View sessions
sqlite3 ./.termai/termai.db ".schema"                    # View schema

# Code generation (if modifying SQL)
sqlc generate                # Generate type-safe Go code from SQL
```

## 🚨 Troubleshooting

### **Common Issues**

#### **1. Import Path Errors**
```bash
# Error: no required module provides package
go mod tidy
go clean -modcache
go mod download
```

#### **2. Missing API Key**
```bash
# Error: provider is not enabled
# Solution: Set API key in .termai.yaml or environment variable
export ANTHROPIC_API_KEY="your-api-key-here"
```

#### **3. Database Permission Issues**
```bash
# Error: failed to create data directory
# Solution: Check permissions
chmod 755 .termai/
```

#### **4. Build Issues**
```bash
# Clean and rebuild
go clean
go mod tidy
go build .
```

### **Getting Help**

#### **Check Configuration**
```bash
./termai --debug    # Enable debug logging to see configuration
```

#### **Verify Setup**
```bash
go version          # Check Go version (need 1.23.5+)
ls -la .termai.yaml # Check config file exists
echo $ANTHROPIC_API_KEY  # Check environment variables
```

#### **Test Database**
```bash
# Check if database is created properly
ls -la .termai/
sqlite3 ./.termai/termai.db ".tables"
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

## Database System Architecture

TermAI uses a sophisticated **SQLite-based persistence layer** with modern tooling for type safety and schema management.

### Core Database Technologies

- **SQLite**: Embedded database for zero-config persistence
- **SQLC**: Type-safe Go code generation from SQL queries
- **golang-migrate**: Schema versioning and migration management
- **WAL Mode**: Write-Ahead Logging for better concurrent performance

### Database Schema Design

#### **Sessions Table**
```sql
CREATE TABLE sessions (
    id TEXT PRIMARY KEY,                    -- UUID session identifier
    parent_session_id TEXT,                 -- Optional parent for session threading
    title TEXT NOT NULL,                    -- Human-readable session title
    message_count INTEGER DEFAULT 0,        -- Cached count with constraints
    prompt_tokens INTEGER DEFAULT 0,        -- LLM usage tracking
    completion_tokens INTEGER DEFAULT 0,    -- LLM usage tracking
    cost REAL DEFAULT 0.0,                 -- Calculated cost tracking
    updated_at INTEGER NOT NULL,           -- Unix timestamp (ms)
    created_at INTEGER NOT NULL            -- Unix timestamp (ms)
);
```

#### **Messages Table**
```sql
CREATE TABLE messages (
    id TEXT PRIMARY KEY,                    -- UUID message identifier
    session_id TEXT NOT NULL,              -- Foreign key to sessions
    role TEXT NOT NULL,                    -- 'user', 'assistant', 'system'
    content TEXT NOT NULL,                 -- Main message content
    thinking TEXT DEFAULT '',              -- AI reasoning/thought process
    finished BOOLEAN DEFAULT 0,            -- Completion status
    tool_calls TEXT,                       -- JSON: LLM tool invocations
    tool_results TEXT,                     -- JSON: Tool execution results
    created_at INTEGER NOT NULL,          -- Unix timestamp (ms)
    updated_at INTEGER NOT NULL,          -- Unix timestamp (ms)
    FOREIGN KEY (session_id) REFERENCES sessions (id) ON DELETE CASCADE
);
```

### Database Triggers & Automation

#### **Automatic Timestamp Updates**
```sql
-- Auto-update timestamps on row changes
CREATE TRIGGER update_sessions_updated_at
AFTER UPDATE ON sessions
BEGIN
    UPDATE sessions SET updated_at = strftime('%s', 'now')
    WHERE id = new.id;
END;
```

#### **Message Count Synchronization**
```sql
-- Maintain accurate message counts automatically
CREATE TRIGGER update_session_message_count_on_insert
AFTER INSERT ON messages
BEGIN
    UPDATE sessions SET message_count = message_count + 1
    WHERE id = new.session_id;
END;
```

### SQLC Code Generation

**Configuration** (`sqlc.yaml`):
```yaml
version: "2"
sql:
  - engine: "sqlite"
    schema: "internal/db/migrations"
    queries: "internal/db/sql"
    gen:
      go:
        package: "db"
        out: "internal/db"
        emit_json_tags: true
        emit_prepared_queries: true
        emit_interface: true
```

**Generated Go Types**:
- Type-safe structs for `Session` and `Message`
- Prepared statement methods for all CRUD operations
- Context-aware query execution
- Transaction support with `WithTx()`

### Migration System

**Directory Structure**:
```
internal/db/migrations/
├── 000001_initial.up.sql    # Schema creation
└── 000001_initial.down.sql  # Schema rollback
```

**Automatic Migration**:
- Embedded migrations using Go embed
- Runs automatically on application startup
- WAL mode and foreign key constraints enabled
- Performance optimizations (page size, cache size)

### Database Operations

#### **Session Management**
```sql
-- Create new session
INSERT INTO sessions (id, title, ...) VALUES (?, ?, ...)

-- List all root sessions (excluding child sessions)
SELECT * FROM sessions WHERE parent_session_id IS NULL 
ORDER BY created_at DESC

-- Update session with usage stats
UPDATE sessions SET prompt_tokens = ?, completion_tokens = ?, cost = ?
WHERE id = ?
```

#### **Message Threading**
```sql
-- Get conversation history for a session
SELECT * FROM messages WHERE session_id = ? ORDER BY created_at ASC

-- Create new message with tool support
INSERT INTO messages (id, session_id, role, content, tool_calls, tool_results)
VALUES (?, ?, ?, ?, ?, ?)

-- Update message during streaming/processing
UPDATE messages SET content = ?, thinking = ?, finished = ?
WHERE id = ?
```

### Performance Features

#### **Database Optimizations**
- **Foreign Keys**: Enabled for referential integrity
- **WAL Mode**: Better concurrent read/write performance
- **Indexes**: `idx_messages_session_id` for fast message queries
- **Page Size**: Optimized 4KB pages
- **Cache Size**: 8MB memory cache

#### **Type Safety Benefits**
- **Compile-time Safety**: SQLC generates type-safe Go code
- **No ORM Overhead**: Direct SQL with zero runtime cost
- **Prepared Statements**: Automatic SQL injection protection
- **Context Support**: Proper cancellation and timeout handling

### Data Flow Architecture

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   TUI Layer     │───▶│  Service Layer   │───▶│  Database Layer │
│                 │    │                  │    │                 │
│ • Session List  │    │ • Session Svc    │    │ • SQLite DB     │
│ • Message View  │    │ • Message Svc    │    │ • SQLC Queries  │
│ • Editor        │    │ • LLM Service    │    │ • Migrations    │
└─────────────────┘    └──────────────────┘    └─────────────────┘
                                │
                       ┌────────▼─────────┐
                       │   Pub/Sub Events │
                       │                  │
                       │ • Session Events │
                       │ • Message Events │
                       │ • Real-time UI   │
                       └──────────────────┘
```

### Database Development Workflow

```bash
# Generate type-safe Go code from SQL
sqlc generate

# View database schema
sqlite3 ./data/termai.db ".schema"

# Query sessions and messages
sqlite3 ./data/termai.db "SELECT * FROM sessions;"
sqlite3 ./data/termai.db "SELECT * FROM messages WHERE session_id = 'uuid';"

# Check database integrity
sqlite3 ./data/termai.db "PRAGMA integrity_check;"
```

This database architecture provides **enterprise-grade persistence** with automatic migrations, type safety, and optimized performance for the AI assistant's conversational data.

## LLM Service Flow Architecture

TermAI implements a sophisticated **event-driven LLM processing pipeline** that handles user prompts through multiple service layers with real-time streaming and tool execution.

### Complete User Prompt → AI Response Flow

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                          USER INPUT FLOW                                       │
└─────────────────────────────────────────────────────────────────────────────────┘

1. User Types Message
   ├─ TUI Editor Component (Vim-style editor)
   ├─ Ctrl+Enter → Send() command
   └─ Validation (session selected, no pending messages)

2. Editor → Agent Creation
   ├─ agent.NewCoderAgent(app)
   ├─ Load configuration & model settings
   └─ Initialize provider (OpenAI/Anthropic/Gemini/Groq)

3. Message Persistence 
   ├─ Messages.Create(sessionID, User role, content)
   ├─ Database INSERT with UUID generation
   └─ Pub/Sub Event → UI updates session list

┌─────────────────────────────────────────────────────────────────────────────────┐
│                        LLM PROCESSING PIPELINE                                 │
└─────────────────────────────────────────────────────────────────────────────────┘

4. Context Building
   ├─ Messages.List(sessionID) → Load conversation history
   ├─ System prompt injection (CoderSystemPrompt)
   └─ Tool definitions attachment (bash, view, edit, write, etc.)

5. Provider Streaming
   ├─ provider.StreamResponse(context, messages, tools)
   ├─ Real-time event channel creation
   └─ HTTP/WebSocket to LLM provider

6. Event Processing Loop
   ├─ EventThinkingDelta → Update assistant thinking field
   ├─ EventContentDelta → Stream response content 
   ├─ EventComplete → Finalize message & usage tracking
   └─ Real-time database updates via Messages.Update()

┌─────────────────────────────────────────────────────────────────────────────────┐
│                         TOOL EXECUTION SYSTEM                                  │
└─────────────────────────────────────────────────────────────────────────────────┘

7. Tool Call Detection
   ├─ Parse tool_calls from LLM response
   ├─ Concurrent execution via goroutines
   └─ Permission system integration

8. Tool Execution (Parallel)
   ├─ Bash Tool → Shell command execution
   ├─ File Tools → read/write/edit operations  
   ├─ Search Tools → glob/grep operations
   └─ Agent Tool → Recursive sub-agent calls

9. Tool Results
   ├─ Collect all tool responses
   ├─ Messages.Create(Tool role, tool_results)
   └─ Continue conversation loop if more tools needed

┌─────────────────────────────────────────────────────────────────────────────────┐
│                       REAL-TIME UI UPDATES                                     │
└─────────────────────────────────────────────────────────────────────────────────┘

10. Pub/Sub Event Distribution
    ├─ Message events → Messages component updates
    ├─ Session events → Sessions list refresh
    └─ Logger events → Status notifications

11. TUI Rendering
    ├─ Bubble Tea event handling
    ├─ Real-time markdown rendering
    └─ Streaming text display in message view
```

### Detailed Service Layer Architecture

#### **1. TUI Layer** (`internal/tui/`)
```go
// Editor sends user input
func (m *editorCmp) Send() tea.Cmd {
    content := strings.Join(m.editor.GetBuffer().Lines(), "\n")
    agent, _ := agent.NewCoderAgent(m.app)
    go agent.Generate(m.sessionID, content)  // Async processing
}
```

#### **2. Agent Layer** (`internal/llm/agent/`)
```go
// Main generation flow with tool execution loop
func (c *agent) generate(sessionID, content string) error {
    // 1. Create user message
    userMsg := c.Messages.Create(sessionID, User, content)
    
    for {
        // 2. Stream LLM response
        eventChan := c.agent.StreamResponse(context, messages, tools)
        assistantMsg := c.Messages.Create(sessionID, Assistant, "")
        
        // 3. Process streaming events
        for event := range eventChan {
            c.processEvent(sessionID, &assistantMsg, event)
        }
        
        // 4. Execute tools if present
        if len(assistantMsg.ToolCalls) > 0 {
            toolResults := c.ExecuteTools(context, assistantMsg.ToolCalls)
            c.Messages.Create(sessionID, Tool, toolResults)
            continue // Loop for more LLM responses
        }
        break // No more tools needed
    }
}
```

#### **3. Provider Layer** (`internal/llm/provider/`)
```go
// Streaming provider interface
type Provider interface {
    StreamResponse(ctx context.Context, messages []Message, tools []Tool) 
        (<-chan ProviderEvent, error)
}

// Event types for real-time updates
const (
    EventThinkingDelta  // AI reasoning process
    EventContentDelta   // Response content chunks  
    EventComplete       // Final response with usage
    EventError          // Error handling
)
```

#### **4. Tool System** (`internal/llm/tools/`)
```go
// Concurrent tool execution
func (c *agent) ExecuteTools(toolCalls []ToolCall, tools []BaseTool) []ToolResult {
    var wg sync.WaitGroup
    results := make([]ToolResult, len(toolCalls))
    
    for i, call := range toolCalls {
        wg.Add(1)
        go func(index int, toolCall ToolCall) {
            defer wg.Done()
            result := tool.Run(context, toolCall)  // Parallel execution
            results[index] = result
        }(i, call)
    }
    wg.Wait()
    return results
}
```

### Advanced Features

#### **Real-time Streaming**
- **Thinking Delta**: Shows AI reasoning process in real-time
- **Content Delta**: Streams response text as it's generated
- **Tool Execution**: Live updates during command execution
- **Usage Tracking**: Token counts and cost calculation

#### **Tool Integration**
- **Bash Tool**: Persistent shell with command execution
- **File Tools**: Read, write, edit operations with permission system
- **Search Tools**: Glob pattern matching and grep functionality
- **Agent Tool**: Recursive sub-agent delegation for complex tasks

#### **State Management**
- **Session Persistence**: All conversations stored in SQLite
- **Message Threading**: Complete conversation history with roles
- **Tool Results**: Structured storage of tool execution results
- **Usage Analytics**: Token consumption and cost tracking per session

#### **Error Handling & Safety**
- **Permission System**: User confirmation for dangerous operations
- **Command Filtering**: Banned commands list for security
- **Timeout Management**: Configurable timeouts for tool execution
- **Error Recovery**: Graceful error handling with user feedback



This architecture provides **production-ready LLM integration** with streaming responses, tool execution, persistent conversations, and real-time user interface updates.

