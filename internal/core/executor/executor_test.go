package executor

import (
	"redis-repo/internal/constant"
	"redis-repo/internal/data_structure"
	"strings"
	"testing"
	"time"
)

func resetGlobalDict() {
	dictStore = data_structure.NewDict()
}

func assertResponse(t *testing.T, got []byte, expected string) {
	gotStr := string(got)
	if gotStr != expected {
		t.Errorf("Expected response %q, got %q", expected, gotStr)
	}
}

// Test PING command
func TestExecutePing(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected string
	}{
		{
			name:     "PING with no arguments",
			args:     []string{},
			expected: "+PONG\r\n",
		},
		{
			name:     "PING with message",
			args:     []string{"Hello"},
			expected: "$5\r\nHello\r\n",
		},
		{
			name:     "PING with multiple arguments (should error)",
			args:     []string{"Hello", "World"},
			expected: "-wrong number of arguments for 'ping' command\r\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := executePing(tt.args)
			assertResponse(t, result, tt.expected)
		})
	}
}

// Test GET command
func TestExecuteGet(t *testing.T) {
	resetGlobalDict()

	tests := []struct {
		name     string
		setup    func()
		args     []string
		expected string
	}{
		{
			name: "GET existing key",
			setup: func() {
				dictStore.Set("testkey", "testvalue", 0)
			},
			args:     []string{"testkey"},
			expected: "$9\r\ntestvalue\r\n",
		},
		{
			name: "GET non-existing key",
			setup: func() {
				// No setup - empty dict
			},
			args:     []string{"nonexistent"},
			expected: constant.RespNil,
		},
		{
			name: "GET with wrong number of arguments",
			setup: func() {
				// No setup needed
			},
			args:     []string{},
			expected: "$47\r\nERR wrong number of arguments for 'GET' command\r\n",
		},
		{
			name: "GET with empty key",
			setup: func() {
				// No setup needed
			},
			args:     []string{""},
			expected: "$13\r\nERR empty key\r\n",
		},
		{
			name: "GET expired key",
			setup: func() {
				// Set key with immediate expiry
				dictStore.Set("expired", "value", uint64(time.Now().UnixMilli()-1000))
			},
			args:     []string{"expired"},
			expected: constant.RespNil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetGlobalDict()
			tt.setup()
			result := executeGet(tt.args)
			assertResponse(t, result, tt.expected)
		})
	}
}

// Test SET command
func TestExecuteSet(t *testing.T) {
	resetGlobalDict()

	tests := []struct {
		name     string
		args     []string
		expected string
	}{
		{
			name:     "SET basic key-value",
			args:     []string{"key", "value"},
			expected: constant.RespOk,
		},
		{
			name:     "SET with EX (seconds)",
			args:     []string{"key", "value", "EX", "60"},
			expected: constant.RespOk,
		},
		{
			name:     "SET with PX (milliseconds)",
			args:     []string{"key", "value", "PX", "60000"},
			expected: constant.RespOk,
		},
		{
			name:     "SET with EXAT (timestamp)",
			args:     []string{"key", "value", "EXAT", "9999999999"},
			expected: constant.RespOk,
		},
		{
			name:     "SET with PXAT (milliseconds timestamp)",
			args:     []string{"key", "value", "PXAT", "9999999999000"},
			expected: constant.RespOk,
		},
		{
			name:     "SET with wrong number of arguments",
			args:     []string{"key"},
			expected: "$47\r\nERR wrong number of arguments for 'SET' command\r\n",
		},
		{
			name:     "SET with empty key",
			args:     []string{"", "value"},
			expected: "$13\r\nERR empty key\r\n",
		},
		{
			name:     "SET with invalid EX value",
			args:     []string{"key", "value", "EX", "invalid"},
			expected: "$16\r\nERR invalid time\r\n",
		},
		{
			name:     "SET with negative EX value",
			args:     []string{"key", "value", "EX", "-1"},
			expected: "$16\r\nERR invalid time\r\n",
		},
		{
			name:     "SET with invalid expiry type",
			args:     []string{"key", "value", "INVALID", "60"},
			expected: "-invalid type of expiry time\r\n",
		},
		{
			name:     "SET with EXAT in the past",
			args:     []string{"key", "value", "EXAT", "1"},
			expected: "$16\r\nERR invalid time\r\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetGlobalDict()
			result := executeSet(tt.args)
			assertResponse(t, result, tt.expected)
		})
	}
}

// Test TTL command
func TestExecuteTTL(t *testing.T) {
	resetGlobalDict()

	tests := []struct {
		name     string
		setup    func()
		args     []string
		expected string
	}{
		{
			name: "TTL for non-existing key",
			setup: func() {
				// No setup - empty dict
			},
			args:     []string{"nonexistent"},
			expected: constant.TtlKeyNotExist,
		},
		{
			name: "TTL for key without expiry",
			setup: func() {
				dictStore.Set("key", "value", 0)
			},
			args:     []string{"key"},
			expected: constant.TtlKeyExistNoExpire,
		},
		{
			name: "TTL for key with future expiry",
			setup: func() {
				futureTime := uint64(time.Now().UnixMilli() + 60000) // 60 seconds from now
				dictStore.Set("key", "value", futureTime)
			},
			args:     []string{"key"},
			expected: ":", // Should return a positive integer (seconds)
		},
		{
			name: "TTL for expired key",
			setup: func() {
				pastTime := uint64(time.Now().UnixMilli() - 1000) // 1 second ago
				dictStore.Set("expired", "value", pastTime)
			},
			args:     []string{"expired"},
			expected: constant.TtlKeyNotExist,
		},
		{
			name: "TTL with wrong number of arguments",
			setup: func() {
				// No setup needed
			},
			args:     []string{},
			expected: "$47\r\nERR wrong number of arguments for 'TTL' command\r\n",
		},
		{
			name: "TTL with empty key",
			setup: func() {
				// No setup needed
			},
			args:     []string{""},
			expected: "$13\r\nERR empty key\r\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetGlobalDict()
			tt.setup()
			result := executeTTL(tt.args)

			// For TTL with future expiry, just check it's a positive integer
			if tt.name == "TTL for key with future expiry" {
				resultStr := string(result)
				if !strings.HasPrefix(resultStr, ":") || strings.Contains(resultStr, "-") {
					t.Errorf("Expected positive integer response, got %q", resultStr)
				}
			} else {
				assertResponse(t, result, tt.expected)
			}
		})
	}
}

// Test DEL command
func TestExecuteDel(t *testing.T) {
	resetGlobalDict()

	tests := []struct {
		name     string
		setup    func()
		args     []string
		expected string
	}{
		{
			name: "DEL existing key",
			setup: func() {
				dictStore.Set("key1", "value1", 0)
			},
			args:     []string{"key1"},
			expected: ":1\r\n",
		},
		{
			name: "DEL non-existing key",
			setup: func() {
				// No setup - empty dict
			},
			args:     []string{"nonexistent"},
			expected: ":0\r\n",
		},
		{
			name: "DEL multiple keys - some exist",
			setup: func() {
				dictStore.Set("key1", "value1", 0)
				dictStore.Set("key2", "value2", 0)
			},
			args:     []string{"key1", "key2", "nonexistent"},
			expected: ":2\r\n",
		},
		{
			name: "DEL multiple keys - none exist",
			setup: func() {
				// No setup - empty dict
			},
			args:     []string{"nonexistent1", "nonexistent2"},
			expected: ":0\r\n",
		},
		{
			name: "DEL with no arguments",
			setup: func() {
				// No setup needed
			},
			args:     []string{},
			expected: ":0\r\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetGlobalDict()
			tt.setup()
			result := executeDel(tt.args)
			assertResponse(t, result, tt.expected)
		})
	}
}

// Integration tests
func TestCommandIntegration(t *testing.T) {
	resetGlobalDict()

	t.Run("SET-GET-TTL-DEL workflow", func(t *testing.T) {
		// SET a key with expiry
		setResult := executeSet([]string{"testkey", "testvalue", "EX", "60"})
		assertResponse(t, setResult, constant.RespOk)

		// GET the key
		getResult := executeGet([]string{"testkey"})
		assertResponse(t, getResult, "$9\r\ntestvalue\r\n")

		// Check TTL
		ttlResult := executeTTL([]string{"testkey"})
		ttlStr := string(ttlResult)
		if !strings.HasPrefix(ttlStr, ":") || strings.Contains(ttlStr, "-") {
			t.Errorf("Expected positive TTL, got %q", ttlStr)
		}

		// DELETE the key
		delResult := executeDel([]string{"testkey"})
		assertResponse(t, delResult, ":1\r\n")

		// GET should return nil after deletion
		getResultAfterDel := executeGet([]string{"testkey"})
		assertResponse(t, getResultAfterDel, constant.RespNil)

		// TTL should return -2 after deletion
		ttlResultAfterDel := executeTTL([]string{"testkey"})
		assertResponse(t, ttlResultAfterDel, constant.TtlKeyNotExist)
	})

	t.Run("Expired key cleanup", func(t *testing.T) {
		// Set a key with immediate expiry
		immediateExpiry := uint64(time.Now().UnixMilli() - 1000)
		dictStore.Set("expired", "value", immediateExpiry)

		// GET should return nil (key should be cleaned up)
		getResult := executeGet([]string{"expired"})
		assertResponse(t, getResult, constant.RespNil)

		// TTL should return -2 (key not found)
		ttlResult := executeTTL([]string{"expired"})
		assertResponse(t, ttlResult, constant.TtlKeyNotExist)
	})
}
