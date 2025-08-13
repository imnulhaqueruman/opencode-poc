
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
   cd /home/user/development/opencode-poc
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

### Performance Optimizations

```bash
# Database optimizations
- WAL mode for concurrent access
- Prepared statements via SQLC  
- Indexed message queries
- Automatic timestamp triggers

# Memory management
- Streaming responses (no buffering)
- Concurrent tool execution
- Connection pooling
- Event-driven architecture

# User experience
- Real-time UI updates via Pub/Sub
- Non-blocking async processing  
- Progress indicators during long operations
- Vim-style editor with familiar keybindings
```

This architecture provides **production-ready LLM integration** with streaming responses, tool execution, persistent conversations, and real-time user interface updates.



[watch](./resources/puku-test-2.gif)