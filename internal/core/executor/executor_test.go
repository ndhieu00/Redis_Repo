package executor

import (
	"redis-repo/internal/constant"
	"redis-repo/internal/data_structure"
	"strings"
	"testing"
	"time"
)

func resetGlobalDict() {
	dict = data_structure.NewDict()
}

func resetGlobalSetStore() {
	setStore = make(map[string]data_structure.Set)
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
			expected: "-ERR wrong number of arguments for 'PING' command\r\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cmdPING(tt.args)
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
				dict.Set("testkey", "testvalue", 0)
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
			expected: "-ERR wrong number of arguments for 'GET' command\r\n",
		},
		{
			name: "GET with empty key",
			setup: func() {
				// No setup needed
			},
			args:     []string{""},
			expected: "-ERR empty key\r\n",
		},
		{
			name: "GET expired key",
			setup: func() {
				// Set key with immediate expiry
				dict.Set("expired", "value", uint64(time.Now().UnixMilli()-1000))
			},
			args:     []string{"expired"},
			expected: constant.RespNil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetGlobalDict()
			tt.setup()
			result := cmdGET(tt.args)
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
			expected: "-ERR wrong number of arguments for 'SET' command\r\n",
		},
		{
			name:     "SET with empty key",
			args:     []string{"", "value"},
			expected: "-ERR empty key\r\n",
		},
		{
			name:     "SET with invalid EX value",
			args:     []string{"key", "value", "EX", "invalid"},
			expected: "-ERR invalid time\r\n",
		},
		{
			name:     "SET with negative EX value",
			args:     []string{"key", "value", "EX", "-1"},
			expected: "-ERR invalid time\r\n",
		},
		{
			name:     "SET with invalid expiry type",
			args:     []string{"key", "value", "INVALID", "60"},
			expected: "-ERR invalid type of expiry time\r\n",
		},
		{
			name:     "SET with EXAT in the past",
			args:     []string{"key", "value", "EXAT", "1"},
			expected: "-ERR invalid time\r\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetGlobalDict()
			result := cmdSET(tt.args)
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
				dict.Set("key", "value", 0)
			},
			args:     []string{"key"},
			expected: constant.TtlKeyExistNoExpire,
		},
		{
			name: "TTL for key with future expiry",
			setup: func() {
				futureTime := uint64(time.Now().UnixMilli() + 60000) // 60 seconds from now
				dict.Set("key", "value", futureTime)
			},
			args:     []string{"key"},
			expected: ":", // Should return a positive integer (seconds)
		},
		{
			name: "TTL for expired key",
			setup: func() {
				pastTime := uint64(time.Now().UnixMilli() - 1000) // 1 second ago
				dict.Set("expired", "value", pastTime)
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
			expected: "-ERR wrong number of arguments for 'TTL' command\r\n",
		},
		{
			name: "TTL with empty key",
			setup: func() {
				// No setup needed
			},
			args:     []string{""},
			expected: "-ERR empty key\r\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetGlobalDict()
			tt.setup()
			result := cmdTTL(tt.args)

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
				dict.Set("key1", "value1", 0)
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
				dict.Set("key1", "value1", 0)
				dict.Set("key2", "value2", 0)
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
			result := cmdDEL(tt.args)
			assertResponse(t, result, tt.expected)
		})
	}
}

// Integration tests
func TestCommandIntegration(t *testing.T) {
	resetGlobalDict()

	t.Run("SET-GET-TTL-DEL workflow", func(t *testing.T) {
		// SET a key with expiry
		setResult := cmdSET([]string{"testkey", "testvalue", "EX", "60"})
		assertResponse(t, setResult, constant.RespOk)

		// GET the key
		getResult := cmdGET([]string{"testkey"})
		assertResponse(t, getResult, "$9\r\ntestvalue\r\n")

		// Check TTL
		ttlResult := cmdTTL([]string{"testkey"})
		ttlStr := string(ttlResult)
		if !strings.HasPrefix(ttlStr, ":") || strings.Contains(ttlStr, "-") {
			t.Errorf("Expected positive TTL, got %q", ttlStr)
		}

		// DELETE the key
		delResult := cmdDEL([]string{"testkey"})
		assertResponse(t, delResult, ":1\r\n")

		// GET should return nil after deletion
		getResultAfterDel := cmdGET([]string{"testkey"})
		assertResponse(t, getResultAfterDel, constant.RespNil)

		// TTL should return -2 after deletion
		ttlResultAfterDel := cmdTTL([]string{"testkey"})
		assertResponse(t, ttlResultAfterDel, constant.TtlKeyNotExist)
	})

	t.Run("Expired key cleanup", func(t *testing.T) {
		// Set a key with immediate expiry
		immediateExpiry := uint64(time.Now().UnixMilli() - 1000)
		dict.Set("expired", "value", immediateExpiry)

		// GET should return nil (key should be cleaned up)
		getResult := cmdGET([]string{"expired"})
		assertResponse(t, getResult, constant.RespNil)

		// TTL should return -2 (key not found)
		ttlResult := cmdTTL([]string{"expired"})
		assertResponse(t, ttlResult, constant.TtlKeyNotExist)
	})
}

// Test SADD command
func TestExecuteSadd(t *testing.T) {
	resetGlobalSetStore()

	tests := []struct {
		name     string
		args     []string
		expected string
	}{
		{
			name:     "SADD new set with single member",
			args:     []string{"myset", "member1"},
			expected: ":1\r\n",
		},
		{
			name:     "SADD new set with multiple members",
			args:     []string{"myset", "member1", "member2", "member3"},
			expected: ":3\r\n",
		},
		{
			name:     "SADD existing set with new members",
			args:     []string{"myset", "member4", "member5"},
			expected: ":2\r\n",
		},
		{
			name:     "SADD existing set with duplicate members",
			args:     []string{"myset", "member1", "member6"},
			expected: ":1\r\n", // Only member6 is new
		},
		{
			name:     "SADD with wrong number of arguments",
			args:     []string{"myset"},
			expected: "-ERR wrong number of arguments for 'SADD' command\r\n",
		},
		{
			name:     "SADD with no arguments",
			args:     []string{},
			expected: "-ERR wrong number of arguments for 'SADD' command\r\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetGlobalSetStore()
			// Pre-populate set for some tests
			if tt.name == "SADD existing set with new members" || tt.name == "SADD existing set with duplicate members" {
				setStore["myset"] = data_structure.NewSet([]string{"member1", "member2", "member3"})
			}
			result := cmdSADD(tt.args)
			assertResponse(t, result, tt.expected)
		})
	}
}

// Test SMEMBERS command
func TestExecuteSmembers(t *testing.T) {
	resetGlobalSetStore()

	tests := []struct {
		name     string
		setup    func()
		args     []string
		expected string
	}{
		{
			name: "SMEMBERS existing set",
			setup: func() {
				setStore["myset"] = data_structure.NewSet([]string{"member1", "member2", "member3"})
			},
			args:     []string{"myset"},
			expected: "*3\r\n", // Just check array length, order is not guaranteed
		},
		{
			name: "SMEMBERS empty set",
			setup: func() {
				setStore["myset"] = data_structure.NewSet([]string{})
			},
			args:     []string{"myset"},
			expected: "*0\r\n",
		},
		{
			name: "SMEMBERS non-existing set",
			setup: func() {
				// No setup - empty setStore
			},
			args:     []string{"nonexistent"},
			expected: "*0\r\n",
		},
		{
			name: "SMEMBERS with wrong number of arguments",
			setup: func() {
				// No setup needed
			},
			args:     []string{},
			expected: "-ERR wrong number of arguments for 'SMEMBERS' command\r\n",
		},
		{
			name: "SMEMBERS with multiple arguments",
			setup: func() {
				// No setup needed
			},
			args:     []string{"myset", "extra"},
			expected: "-ERR wrong number of arguments for 'SMEMBERS' command\r\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetGlobalSetStore()
			tt.setup()
			result := cmdSMEMBERS(tt.args)

			// For existing set test, just check array length since order is not guaranteed
			if tt.name == "SMEMBERS existing set" {
				resultStr := string(result)
				if !strings.HasPrefix(resultStr, "*3\r\n") {
					t.Errorf("Expected array with 3 elements, got %q", resultStr)
				}
			} else {
				assertResponse(t, result, tt.expected)
			}
		})
	}
}

// Test SISMEMBER command
func TestExecuteSismember(t *testing.T) {
	resetGlobalSetStore()

	tests := []struct {
		name     string
		setup    func()
		args     []string
		expected string
	}{
		{
			name: "SISMEMBER existing member",
			setup: func() {
				setStore["myset"] = data_structure.NewSet([]string{"member1", "member2", "member3"})
			},
			args:     []string{"myset", "member1"},
			expected: "*1\r\n:1\r\n",
		},
		{
			name: "SISMEMBER non-existing member",
			setup: func() {
				setStore["myset"] = data_structure.NewSet([]string{"member1", "member2", "member3"})
			},
			args:     []string{"myset", "member4"},
			expected: "*1\r\n:0\r\n",
		},
		{
			name: "SISMEMBER non-existing set",
			setup: func() {
				// No setup - empty setStore
			},
			args:     []string{"nonexistent", "member1"},
			expected: "*1\r\n:0\r\n",
		},
		{
			name: "SISMEMBER multiple members - mixed results",
			setup: func() {
				setStore["myset"] = data_structure.NewSet([]string{"member1", "member2", "member3"})
			},
			args:     []string{"myset", "member1", "member4", "member2"},
			expected: "*3\r\n:1\r\n:0\r\n:1\r\n",
		},
		{
			name: "SISMEMBER with wrong number of arguments",
			setup: func() {
				// No setup needed
			},
			args:     []string{"myset"},
			expected: "-ERR wrong number of arguments for 'SMISMEMBER' command\r\n",
		},
		{
			name: "SISMEMBER with no arguments",
			setup: func() {
				// No setup needed
			},
			args:     []string{},
			expected: "-ERR wrong number of arguments for 'SMISMEMBER' command\r\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetGlobalSetStore()
			tt.setup()
			result := cmdSMISMEMBER(tt.args)
			assertResponse(t, result, tt.expected)
		})
	}
}

// Test SREM command
func TestExecuteSrem(t *testing.T) {
	resetGlobalSetStore()

	tests := []struct {
		name     string
		setup    func()
		args     []string
		expected string
	}{
		{
			name: "SREM existing member",
			setup: func() {
				setStore["myset"] = data_structure.NewSet([]string{"member1", "member2", "member3"})
			},
			args:     []string{"myset", "member1"},
			expected: ":1\r\n",
		},
		{
			name: "SREM non-existing member",
			setup: func() {
				setStore["myset"] = data_structure.NewSet([]string{"member1", "member2", "member3"})
			},
			args:     []string{"myset", "member4"},
			expected: ":0\r\n",
		},
		{
			name: "SREM non-existing set",
			setup: func() {
				// No setup - empty setStore
			},
			args:     []string{"nonexistent", "member1"},
			expected: ":0\r\n",
		},
		{
			name: "SREM multiple members - some exist",
			setup: func() {
				setStore["myset"] = data_structure.NewSet([]string{"member1", "member2", "member3"})
			},
			args:     []string{"myset", "member1", "member4", "member2"},
			expected: ":2\r\n", // Only member1 and member2 were removed
		},
		{
			name: "SREM multiple members - none exist",
			setup: func() {
				setStore["myset"] = data_structure.NewSet([]string{"member1", "member2", "member3"})
			},
			args:     []string{"myset", "member4", "member5"},
			expected: ":0\r\n",
		},
		{
			name: "SREM with wrong number of arguments",
			setup: func() {
				// No setup needed
			},
			args:     []string{"myset"},
			expected: "-ERR wrong number of arguments for 'SREM' command\r\n",
		},
		{
			name: "SREM with no arguments",
			setup: func() {
				// No setup needed
			},
			args:     []string{},
			expected: "-ERR wrong number of arguments for 'SREM' command\r\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetGlobalSetStore()
			tt.setup()
			result := cmdSREM(tt.args)
			assertResponse(t, result, tt.expected)
		})
	}
}

// Integration tests for set commands
func TestSetCommandIntegration(t *testing.T) {
	resetGlobalSetStore()

	t.Run("SADD-SMEMBERS-SISMEMBER-SREM workflow", func(t *testing.T) {
		// SADD members to a new set
		saddResult := cmdSADD([]string{"myset", "member1", "member2", "member3"})
		assertResponse(t, saddResult, ":3\r\n")

		// SMEMBERS should return all members
		smembersResult := cmdSMEMBERS([]string{"myset"})
		smembersStr := string(smembersResult)
		if !strings.HasPrefix(smembersStr, "*3\r\n") {
			t.Errorf("Expected array with 3 elements, got %q", smembersStr)
		}

		// SISMEMBER should return 1 for existing members
		sismemberResult := cmdSMISMEMBER([]string{"myset", "member1", "member4"})
		assertResponse(t, sismemberResult, "*2\r\n:1\r\n:0\r\n")

		// SADD more members (some duplicates)
		saddMoreResult := cmdSADD([]string{"myset", "member3", "member4", "member5"})
		assertResponse(t, saddMoreResult, ":2\r\n") // Only member4 and member5 are new

		// SMEMBERS should now have 5 members
		smembersAfterAddResult := cmdSMEMBERS([]string{"myset"})
		// Note: Order is not guaranteed in sets, so we just check it's an array with 5 elements
		smembersAfterAddStr := string(smembersAfterAddResult)
		if !strings.HasPrefix(smembersAfterAddStr, "*5\r\n") {
			t.Errorf("Expected array with 5 elements, got %q", smembersAfterAddStr)
		}

		// SREM some members
		sremResult := cmdSREM([]string{"myset", "member1", "member6"})
		assertResponse(t, sremResult, ":1\r\n") // Only member1 was removed

		// Final SMEMBERS should have 4 members
		finalSmembersResult := cmdSMEMBERS([]string{"myset"})
		finalSmembersStr := string(finalSmembersResult)
		if !strings.HasPrefix(finalSmembersStr, "*4\r\n") {
			t.Errorf("Expected array with 4 elements, got %q", finalSmembersStr)
		}
	})

	t.Run("Empty set operations", func(t *testing.T) {
		// SMEMBERS on non-existing set
		smembersResult := cmdSMEMBERS([]string{"empty"})
		assertResponse(t, smembersResult, "*0\r\n")

		// SISMEMBER on non-existing set
		sismemberResult := cmdSMISMEMBER([]string{"empty", "member1"})
		assertResponse(t, sismemberResult, "*1\r\n:0\r\n")

		// SREM on non-existing set
		sremResult := cmdSREM([]string{"empty", "member1"})
		assertResponse(t, sremResult, ":0\r\n")
	})
}

// Test SCARD command
func TestExecuteScard(t *testing.T) {
	resetGlobalSetStore()

	tests := []struct {
		name     string
		setup    func()
		args     []string
		expected string
	}{
		{
			name: "SCARD existing set with members",
			setup: func() {
				setStore["myset"] = data_structure.NewSet([]string{"member1", "member2", "member3"})
			},
			args:     []string{"myset"},
			expected: ":3\r\n",
		},
		{
			name: "SCARD empty set",
			setup: func() {
				setStore["myset"] = data_structure.NewSet([]string{})
			},
			args:     []string{"myset"},
			expected: ":0\r\n",
		},
		{
			name: "SCARD non-existing set",
			setup: func() {
				// No setup - empty setStore
			},
			args:     []string{"nonexistent"},
			expected: ":0\r\n",
		},
		{
			name: "SCARD with wrong number of arguments",
			setup: func() {
				// No setup needed
			},
			args:     []string{},
			expected: "-ERR wrong number of arguments for 'SCARD' command\r\n",
		},
		{
			name: "SCARD with multiple arguments",
			setup: func() {
				// No setup needed
			},
			args:     []string{"myset", "extra"},
			expected: "-ERR wrong number of arguments for 'SCARD' command\r\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetGlobalSetStore()
			tt.setup()
			result := cmdSCARD(tt.args)
			assertResponse(t, result, tt.expected)
		})
	}
}

// Test SINTER command
func TestExecuteSinter(t *testing.T) {
	resetGlobalSetStore()

	tests := []struct {
		name     string
		setup    func()
		args     []string
		expected string
	}{
		{
			name: "SINTER two sets with common members",
			setup: func() {
				setStore["set1"] = data_structure.NewSet([]string{"a", "b", "c"})
				setStore["set2"] = data_structure.NewSet([]string{"b", "c", "d"})
			},
			args:     []string{"set1", "set2"},
			expected: "*2\r\n", // Should have 2 common members (b, c)
		},
		{
			name: "SINTER three sets with common members",
			setup: func() {
				setStore["set1"] = data_structure.NewSet([]string{"a", "b", "c", "d"})
				setStore["set2"] = data_structure.NewSet([]string{"b", "c", "d", "e"})
				setStore["set3"] = data_structure.NewSet([]string{"c", "d", "e", "f"})
			},
			args:     []string{"set1", "set2", "set3"},
			expected: "*2\r\n", // Should have 2 common members (c, d)
		},
		{
			name: "SINTER sets with no common members",
			setup: func() {
				setStore["set1"] = data_structure.NewSet([]string{"a", "b"})
				setStore["set2"] = data_structure.NewSet([]string{"c", "d"})
			},
			args:     []string{"set1", "set2"},
			expected: "*0\r\n",
		},
		{
			name: "SINTER identical sets",
			setup: func() {
				setStore["set1"] = data_structure.NewSet([]string{"a", "b", "c"})
				setStore["set2"] = data_structure.NewSet([]string{"a", "b", "c"})
			},
			args:     []string{"set1", "set2"},
			expected: "*3\r\n", // Should have 3 common members
		},
		{
			name: "SINTER with non-existing set",
			setup: func() {
				setStore["set1"] = data_structure.NewSet([]string{"a", "b", "c"})
				// set2 doesn't exist
			},
			args:     []string{"set1", "nonexistent"},
			expected: "*0\r\n",
		},
		{
			name: "SINTER with empty set",
			setup: func() {
				setStore["set1"] = data_structure.NewSet([]string{"a", "b", "c"})
				setStore["set2"] = data_structure.NewSet([]string{})
			},
			args:     []string{"set1", "set2"},
			expected: "*0\r\n",
		},
		{
			name: "SINTER single set",
			setup: func() {
				setStore["set1"] = data_structure.NewSet([]string{"a", "b", "c"})
			},
			args:     []string{"set1"},
			expected: "*3\r\n", // Should return all members of the single set
		},
		{
			name: "SINTER with no arguments",
			setup: func() {
				// No setup needed
			},
			args:     []string{},
			expected: "-ERR wrong number of arguments for 'SINTER' command\r\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetGlobalSetStore()
			tt.setup()
			result := cmdSINTER(tt.args)

			// For SINTER tests, we check array length since order is not guaranteed
			if strings.HasPrefix(tt.expected, "*") {
				resultStr := string(result)
				if !strings.HasPrefix(resultStr, tt.expected) {
					t.Errorf("Expected %s, got %q", tt.expected, resultStr)
				}
			} else {
				assertResponse(t, result, tt.expected)
			}
		})
	}
}

// Integration tests for SCARD and SINTER commands
func TestScardSinterIntegration(t *testing.T) {
	resetGlobalSetStore()

	t.Run("SCARD-SINTER workflow", func(t *testing.T) {
		// Create sets with some overlap
		saddResult1 := cmdSADD([]string{"set1", "a", "b", "c", "d"})
		assertResponse(t, saddResult1, ":4\r\n")

		saddResult2 := cmdSADD([]string{"set2", "c", "d", "e", "f"})
		assertResponse(t, saddResult2, ":4\r\n")

		saddResult3 := cmdSADD([]string{"set3", "d", "e", "f", "g"})
		assertResponse(t, saddResult3, ":4\r\n")

		// Check cardinality of each set
		scardResult1 := cmdSCARD([]string{"set1"})
		assertResponse(t, scardResult1, ":4\r\n")

		scardResult2 := cmdSCARD([]string{"set2"})
		assertResponse(t, scardResult2, ":4\r\n")

		scardResult3 := cmdSCARD([]string{"set3"})
		assertResponse(t, scardResult3, ":4\r\n")

		// Find intersection of set1 and set2
		sinterResult12 := cmdSINTER([]string{"set1", "set2"})
		sinter12Str := string(sinterResult12)
		if !strings.HasPrefix(sinter12Str, "*2\r\n") {
			t.Errorf("Expected intersection of set1 and set2 to have 2 elements, got %q", sinter12Str)
		}

		// Find intersection of all three sets
		sinterResult123 := cmdSINTER([]string{"set1", "set2", "set3"})
		sinter123Str := string(sinterResult123)
		if !strings.HasPrefix(sinter123Str, "*1\r\n") {
			t.Errorf("Expected intersection of all three sets to have 1 element, got %q", sinter123Str)
		}

		// Remove some elements and check cardinality again
		sremResult := cmdSREM([]string{"set1", "a", "b"})
		assertResponse(t, sremResult, ":2\r\n")

		scardAfterRemoval := cmdSCARD([]string{"set1"})
		assertResponse(t, scardAfterRemoval, ":2\r\n")

		// Check intersection after removal
		sinterAfterRemoval := cmdSINTER([]string{"set1", "set2"})
		sinterAfterStr := string(sinterAfterRemoval)
		if !strings.HasPrefix(sinterAfterStr, "*2\r\n") {
			t.Errorf("Expected intersection after removal to have 2 elements, got %q", sinterAfterStr)
		}
	})

	t.Run("Edge cases", func(t *testing.T) {
		// SCARD on non-existing set
		scardResult := cmdSCARD([]string{"nonexistent"})
		assertResponse(t, scardResult, ":0\r\n")

		// SINTER with non-existing sets
		sinterResult := cmdSINTER([]string{"nonexistent1", "nonexistent2"})
		assertResponse(t, sinterResult, "*0\r\n")

		// SINTER with one existing and one non-existing set
		cmdSADD([]string{"existing", "a", "b"})
		sinterMixedResult := cmdSINTER([]string{"existing", "nonexistent"})
		assertResponse(t, sinterMixedResult, "*0\r\n")
	})
}
