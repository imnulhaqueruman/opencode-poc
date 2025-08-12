# Changelog - TermAI (OpenCode Genesis)

This tracks the evolution of commits we've integrated into this POC.

## Commit 2: "add help" (7844cacb) - March 23, 2025

### ✨ Features Added
- **Help System**: Press `?` to toggle contextual help overlay
- **Status Bar**: Bottom status bar showing version and messages
- **Version Management**: Build-time version detection and display
- **Enhanced Navigation**: Better keybinding system with more options

### 🔧 Technical Changes
- Added `internal/tui/components/core/help.go` - Help overlay component
- Added `internal/tui/components/core/status.go` - Status bar component  
- Added `internal/tui/util/util.go` - Utility functions for messages
- Added `internal/version/version.go` - Version management
- Updated `internal/tui/tui.go` - Integration of help and status systems
- Updated layout pointer receivers for better memory management

### ⌨️ New Keybindings
- `?` - Toggle help overlay
- `Esc` - Close help/current view (changed from "back")
- `Backspace` - Go back to previous page (new)
- `L` - Switch to logs page (unchanged)
- `Ctrl+C`/`q` - Quit application (unchanged)

### 🎨 UI Improvements
- Status bar with version display and help hints
- Error messages displayed in red
- Info messages displayed in green
- Help overlay shows all available keyboard shortcuts
- Dynamic height adjustment when help is toggled

### 📊 Stats
- **Files Changed**: 8 files
- **Lines Added**: +322
- **Lines Removed**: -25
- **Net Change**: +297 lines

---

## Commit 1: "initial" (4b0ea68d) - March 21, 2025

### 🎯 Initial Release
- Basic Go TUI application using Bubble Tea framework
- Multi-pane Bento layout system
- Three pages: REPL, Logs, Init
- Structured logging with pub/sub system
- Event-driven architecture foundation
- Markdown rendering support

### 📊 Stats
- **Files Created**: 28 files
- **Lines Added**: +2,229
- **Total Project Foundation**: Complete TUI framework