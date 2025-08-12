package tools

import (
	"context"
	"encoding/json"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/imnulhaqueruman/opencode-poc/internal/permission"
	"github.com/imnulhaqueruman/opencode-poc/internal/pubsub"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBashTool_Info(t *testing.T) {
	tool := NewBashTool()
	info := tool.Info()

	assert.Equal(t, BashToolName, info.Name)
	assert.NotEmpty(t, info.Description)
	assert.Contains(t, info.Parameters, "command")
	assert.Contains(t, info.Parameters, "timeout")
	assert.Contains(t, info.Required, "command")
}

func TestBashTool_Run(t *testing.T) {
	// Setup a mock permission handler that always allows
	origPermission := permission.Default
	defer func() {
		permission.Default = origPermission
	}()
	permission.Default = newMockPermissionService(true)

	// Save original working directory
	origWd, err := os.Getwd()
	require.NoError(t, err)
	defer func() {
		os.Chdir(origWd)
	}()

	t.Run("executes command successfully", func(t *testing.T) {
		permission.Default = newMockPermissionService(true)
		tool := NewBashTool()
		params := BashParams{
			Command: "echo 'Hello World'",
		}

		paramsJSON, err := json.Marshal(params)
		require.NoError(t, err)

		call := ToolCall{
			Name:  BashToolName,
			Input: string(paramsJSON),
		}

		response, err := tool.Run(context.Background(), call)
		require.NoError(t, err)
		assert.Equal(t, "Hello World\n", response.Content)
	})

	t.Run("handles invalid parameters", func(t *testing.T) {
		permission.Default = newMockPermissionService(true)

		tool := NewBashTool()
		call := ToolCall{
			Name:  BashToolName,
			Input: "invalid json",
		}

		response, err := tool.Run(context.Background(), call)
		require.NoError(t, err)
		assert.Contains(t, response.Content, "invalid parameters")
	})

	t.Run("handles missing command", func(t *testing.T) {
		permission.Default = newMockPermissionService(true)

		tool := NewBashTool()
		params := BashParams{
			Command: "",
		}

		paramsJSON, err := json.Marshal(params)
		require.NoError(t, err)

		call := ToolCall{
			Name:  BashToolName,
			Input: string(paramsJSON),
		}

		response, err := tool.Run(context.Background(), call)
		require.NoError(t, err)
		assert.Contains(t, response.Content, "missing command")
	})

	t.Run("handles banned commands", func(t *testing.T) {
		permission.Default = newMockPermissionService(true)

		tool := NewBashTool()

		for _, bannedCmd := range BannedCommands {
			params := BashParams{
				Command: bannedCmd + " arg1 arg2",
			}

			paramsJSON, err := json.Marshal(params)
			require.NoError(t, err)

			call := ToolCall{
				Name:  BashToolName,
				Input: string(paramsJSON),
			}

			response, err := tool.Run(context.Background(), call)
			require.NoError(t, err)
			assert.Contains(t, response.Content, "not allowed", "Command %s should be blocked", bannedCmd)
		}
	})

	t.Run("handles safe read-only commands without permission check", func(t *testing.T) {
		permission.Default = newMockPermissionService(false)

		tool := NewBashTool()

		// Test with a safe read-only command
		params := BashParams{
			Command: "echo 'test'",
		}

		paramsJSON, err := json.Marshal(params)
		require.NoError(t, err)

		call := ToolCall{
			Name:  BashToolName,
			Input: string(paramsJSON),
		}

		response, err := tool.Run(context.Background(), call)
		require.NoError(t, err)
		assert.Equal(t, "test\n", response.Content)
	})

	t.Run("handles permission denied", func(t *testing.T) {
		permission.Default = newMockPermissionService(false)

		tool := NewBashTool()

		// Test with a command that requires permission
		params := BashParams{
			Command: "mkdir test_dir",
		}

		paramsJSON, err := json.Marshal(params)
		require.NoError(t, err)

		call := ToolCall{
			Name:  BashToolName,
			Input: string(paramsJSON),
		}

		response, err := tool.Run(context.Background(), call)
		require.NoError(t, err)
		assert.Contains(t, response.Content, "permission denied")
	})

	t.Run("handles command timeout", func(t *testing.T) {
		permission.Default = newMockPermissionService(true)
		tool := NewBashTool()
		params := BashParams{
			Command: "sleep 2",
			Timeout: 100, // 100ms timeout
		}

		paramsJSON, err := json.Marshal(params)
		require.NoError(t, err)

		call := ToolCall{
			Name:  BashToolName,
			Input: string(paramsJSON),
		}

		response, err := tool.Run(context.Background(), call)
		require.NoError(t, err)
		assert.Contains(t, response.Content, "aborted")
	})

	t.Run("handles command with stderr output", func(t *testing.T) {
		permission.Default = newMockPermissionService(true)
		tool := NewBashTool()
		params := BashParams{
			Command: "echo 'error message' >&2",
		}

		paramsJSON, err := json.Marshal(params)
		require.NoError(t, err)

		call := ToolCall{
			Name:  BashToolName,
			Input: string(paramsJSON),
		}

		response, err := tool.Run(context.Background(), call)
		require.NoError(t, err)
		assert.Contains(t, response.Content, "error message")
	})

	t.Run("handles command with both stdout and stderr", func(t *testing.T) {
		permission.Default = newMockPermissionService(true)
		tool := NewBashTool()
		params := BashParams{
			Command: "echo 'stdout message' && echo 'stderr message' >&2",
		}

		paramsJSON, err := json.Marshal(params)
		require.NoError(t, err)

		call := ToolCall{
			Name:  BashToolName,
			Input: string(paramsJSON),
		}

		response, err := tool.Run(context.Background(), call)
		require.NoError(t, err)
		assert.Contains(t, response.Content, "stdout message")
		assert.Contains(t, response.Content, "stderr message")
	})

	t.Run("handles context cancellation", func(t *testing.T) {
		permission.Default = newMockPermissionService(true)
		tool := NewBashTool()
		params := BashParams{
			Command: "sleep 5",
		}

		paramsJSON, err := json.Marshal(params)
		require.NoError(t, err)

		call := ToolCall{
			Name:  BashToolName,
			Input: string(paramsJSON),
		}

		ctx, cancel := context.WithCancel(context.Background())

		// Cancel the context after a short delay
		go func() {
			time.Sleep(100 * time.Millisecond)
			cancel()
		}()

		response, err := tool.Run(ctx, call)
		require.NoError(t, err)
		assert.Contains(t, response.Content, "aborted")
	})

	t.Run("respects max timeout", func(t *testing.T) {
		permission.Default = newMockPermissionService(true)
		tool := NewBashTool()
		params := BashParams{
			Command: "echo 'test'",
			Timeout: MaxTimeout + 1000, // Exceeds max timeout
		}

		paramsJSON, err := json.Marshal(params)
		require.NoError(t, err)

		call := ToolCall{
			Name:  BashToolName,
			Input: string(paramsJSON),
		}

		response, err := tool.Run(context.Background(), call)
		require.NoError(t, err)
		assert.Equal(t, "test\n", response.Content)
	})

	t.Run("uses default timeout for zero or negative timeout", func(t *testing.T) {
		permission.Default = newMockPermissionService(true)
		tool := NewBashTool()
		params := BashParams{
			Command: "echo 'test'",
			Timeout: -100, // Negative timeout
		}

		paramsJSON, err := json.Marshal(params)
		require.NoError(t, err)

		call := ToolCall{
			Name:  BashToolName,
			Input: string(paramsJSON),
		}

		response, err := tool.Run(context.Background(), call)
		require.NoError(t, err)
		assert.Equal(t, "test\n", response.Content)
	})
}

func TestTruncateOutput(t *testing.T) {
	t.Run("does not truncate short output", func(t *testing.T) {
		output := "short output"
		result := truncateOutput(output)
		assert.Equal(t, output, result)
	})

	t.Run("truncates long output", func(t *testing.T) {
		// Create a string longer than MaxOutputLength
		longOutput := strings.Repeat("a\n", MaxOutputLength)
		result := truncateOutput(longOutput)

		// Check that the result is shorter than the original
		assert.Less(t, len(result), len(longOutput))

		// Check that the truncation message is included
		assert.Contains(t, result, "lines truncated")

		// Check that we have the beginning and end of the original string
		assert.True(t, strings.HasPrefix(result, "a\n"))
		assert.True(t, strings.HasSuffix(result, "a\n"))
	})
}

func TestCountLines(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected int
	}{
		{
			name:     "empty string",
			input:    "",
			expected: 0,
		},
		{
			name:     "single line",
			input:    "line1",
			expected: 1,
		},
		{
			name:     "multiple lines",
			input:    "line1\nline2\nline3",
			expected: 3,
		},
		{
			name:     "trailing newline",
			input:    "line1\nline2\n",
			expected: 3, // Empty string after last newline counts as a line
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := countLines(tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}

// Mock permission service for testing
type mockPermissionService struct {
	*pubsub.Broker[permission.PermissionRequest]
	allow bool
}

func (m *mockPermissionService) GrantPersistant(permission permission.PermissionRequest) {
	// Not needed for tests
}

func (m *mockPermissionService) Grant(permission permission.PermissionRequest) {
	// Not needed for tests
}

func (m *mockPermissionService) Deny(permission permission.PermissionRequest) {
	// Not needed for tests
}

func (m *mockPermissionService) Request(opts permission.CreatePermissionRequest) bool {
	return m.allow
}

func newMockPermissionService(allow bool) permission.Service {
	return &mockPermissionService{
		Broker: pubsub.NewBroker[permission.PermissionRequest](),
		allow:  allow,
	}
}