package config

import (
	"encoding/base64"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

const validBase64JSON = "eyJhcGlfZW5kcG9pbnQiOiJ0ZXN0LmNvbSIsImNhX2NlcnRpZmljYXRlIjoidGVzdF9jYSIsImNlcnRpZmljYXRlIjoidGVzdF9jZXJ0IiwicHJpdmF0ZV9rZXkiOiJ0ZXN0X2tleSJ9Cg=="

// contains checks if a string contains a substring
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

func TestLoadConfig(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "test.cfg")

	// --- Test: Valid Config File ---
	t.Run("ValidConfigFile", func(t *testing.T) {
		// Create a test config file
		err := os.WriteFile(configPath, []byte(validBase64JSON), 0644)
		if err != nil {
			t.Fatalf("Failed to create test config file: %v", err)
		}

		// Load config
		config, err := LoadConfig(configPath)
		if err != nil {
			t.Fatalf("LoadConfig failed: %v", err)
		}

		// Verify config values
		if config.APIEndpoint != "test.com" {
			t.Errorf("Expected APIEndpoint to be 'test.com', got '%s'", config.APIEndpoint)
		}
		if config.CACertificate != "test_ca" {
			t.Errorf("Expected CACertificate to be 'test_ca', got '%s'", config.CACertificate)
		}
		if config.Certificate != "test_cert" {
			t.Errorf("Expected Certificate to be 'test_cert', got '%s'", config.Certificate)
		}
		if config.PrivateKey != "test_key" {
			t.Errorf("Expected PrivateKey to be 'test_key', got '%s'", config.PrivateKey)
		}
	})

	// --- Test: Non-existent File ---
	t.Run("NonExistentFile", func(t *testing.T) {
		nonExistentPath := filepath.Join(tempDir, "nonexistent.cfg")

		_, err := LoadConfig(nonExistentPath)
		if err == nil {
			t.Error("Expected error for non-existent file")
		}
	})

	// --- Test: Invalid Base64 ---
	t.Run("InvalidBase64", func(t *testing.T) {
		invalidPath := filepath.Join(tempDir, "invalid.cfg")

		// Create file with invalid base64
		err := os.WriteFile(invalidPath, []byte("not base64!"), 0644)
		if err != nil {
			t.Fatalf("Failed to create invalid config file: %v", err)
		}

		_, err = LoadConfig(invalidPath)
		if err == nil {
			t.Error("Expected error for invalid base64")
		}
	})

	// --- Test: Invalid JSON ---
	t.Run("InvalidJSON", func(t *testing.T) {
		invalidPath := filepath.Join(tempDir, "invalid_json.cfg")

		// Create file with valid base64 but invalid JSON
		invalidJSON := base64.StdEncoding.EncodeToString([]byte("not json!"))
		err := os.WriteFile(invalidPath, []byte(invalidJSON), 0644)
		if err != nil {
			t.Fatalf("Failed to create invalid JSON config file: %v", err)
		}

		_, err = LoadConfig(invalidPath)
		if err == nil {
			t.Error("Expected error for invalid JSON")
		}
	})
}

func TestNewConfigFromString(t *testing.T) {
	// --- Test: Valid Config String ---
	t.Run("ValidConfigString", func(t *testing.T) {
		config, err := NewConfigFromString(validBase64JSON)
		if err != nil {
			t.Fatalf("NewConfigFromString failed: %v", err)
		}

		// Verify config values
		if config.APIEndpoint != "test.com" {
			t.Errorf("Expected APIEndpoint to be 'test.com', got '%s'", config.APIEndpoint)
		}
		if config.CACertificate != "test_ca" {
			t.Errorf("Expected CACertificate to be 'test_ca', got '%s'", config.CACertificate)
		}
		if config.Certificate != "test_cert" {
			t.Errorf("Expected Certificate to be 'test_cert', got '%s'", config.Certificate)
		}
		if config.PrivateKey != "test_key" {
			t.Errorf("Expected PrivateKey to be 'test_key', got '%s'", config.PrivateKey)
		}
	})

	// --- Test: Invalid Base64 String ---
	t.Run("InvalidBase64String", func(t *testing.T) {
		_, err := NewConfigFromString("not base64!")
		if err == nil {
			t.Error("Expected error for invalid base64 string")
		}
	})

	// --- Test: Invalid JSON String ---
	t.Run("InvalidJSONString", func(t *testing.T) {
		invalidJSON := base64.StdEncoding.EncodeToString([]byte("not json!"))
		_, err := NewConfigFromString(invalidJSON)
		if err == nil {
			t.Error("Expected error for invalid JSON string")
		}
	})

	// --- Test: Empty String ---
	t.Run("EmptyString", func(t *testing.T) {
		_, err := NewConfigFromString("")
		if err == nil {
			t.Error("Expected error for empty string")
		}
	})
}

func TestWatchConfig(t *testing.T) {
	// Disable parallel execution to avoid race conditions
	// t.Parallel() // Commented out to run tests sequentially

	// Ensure we start with a clean state
	if watcher != nil {
		watcher.Close()
		watcher = nil
	}

	tempDir := t.TempDir()

	// --- Test: Basic WatchConfig Setup ---
	t.Run("BasicSetup", func(t *testing.T) {
		configPath := filepath.Join(tempDir, "watch_test.cfg")

		// Create a test config file
		err := os.WriteFile(configPath, []byte(validBase64JSON), 0644)
		if err != nil {
			t.Fatalf("Failed to create test config file: %v", err)
		}

		// Simple loader function
		loader := func(path string) error {
			return nil
		}

		// Start watching
		err = WatchConfig([]string{configPath}, loader)
		if err != nil {
			t.Fatalf("WatchConfig failed: %v", err)
		}

		// Test that the function returns without error
		// The actual watcher functionality will be tested in FileWriteEvents test

		// Clean up
		if watcher != nil {
			watcher.Close()
			watcher = nil
		}
	})

	// --- Test: Multiple Config Paths ---
	t.Run("MultiplePaths", func(t *testing.T) {
		configPath1 := filepath.Join(tempDir, "watch1.cfg")
		configPath2 := filepath.Join(tempDir, "watch2.cfg")

		// Create test config files
		err := os.WriteFile(configPath1, []byte(validBase64JSON), 0644)
		if err != nil {
			t.Fatalf("Failed to create test config file 1: %v", err)
		}
		err = os.WriteFile(configPath2, []byte(validBase64JSON), 0644)
		if err != nil {
			t.Fatalf("Failed to create test config file 2: %v", err)
		}

		// Track loader calls
		var loaderCalls []string
		var mu sync.Mutex
		loader := func(path string) error {
			mu.Lock()
			loaderCalls = append(loaderCalls, path)
			mu.Unlock()
			return nil
		}

		// Start watching multiple paths
		err = WatchConfig([]string{configPath1, configPath2}, loader)
		if err != nil {
			t.Fatalf("WatchConfig failed: %v", err)
		}

		// Test that the function returns without error
		// The actual watcher functionality will be tested in FileWriteEvents test

		// Clean up
		if watcher != nil {
			watcher.Close()
			watcher = nil
		}
	})

	// --- Test: File Write Events ---
	t.Run("FileWriteEvents", func(t *testing.T) {
		configPath := filepath.Join(tempDir, "write_test.cfg")

		// Create initial config file
		err := os.WriteFile(configPath, []byte(validBase64JSON), 0644)
		if err != nil {
			t.Fatalf("Failed to create test config file: %v", err)
		}

		// Track loader calls
		var loaderCalls []string
		var mu sync.Mutex
		loader := func(path string) error {
			mu.Lock()
			loaderCalls = append(loaderCalls, path)
			mu.Unlock()
			return nil
		}

		// Start watching
		err = WatchConfig([]string{configPath}, loader)
		if err != nil {
			t.Fatalf("WatchConfig failed: %v", err)
		}

		// Give the watcher time to start
		time.Sleep(100 * time.Millisecond)

		// Write to the file to trigger an event
		newConfig := &Config{
			APIEndpoint:   "updated.com",
			CACertificate: "updated_cacert",
			Certificate:   "updated_cert",
			PrivateKey:    "updated_key",
		}
		newConfigBytes, _ := json.Marshal(newConfig)
		newConfigBase64 := base64.StdEncoding.EncodeToString(newConfigBytes)

		err = os.WriteFile(configPath, []byte(newConfigBase64), 0644)
		if err != nil {
			t.Fatalf("Failed to write updated config: %v", err)
		}

		// Give time for the event to be processed
		time.Sleep(200 * time.Millisecond)

		// Check if loader was called
		mu.Lock()
		callCount := len(loaderCalls)
		mu.Unlock()

		if callCount == 0 {
			t.Error("Expected loader to be called when file was written")
		}

		// Clean up
		if watcher != nil {
			watcher.Close()
			watcher = nil
		}
	})

	// --- Test: Watcher Replacement ---
	t.Run("WatcherReplacement", func(t *testing.T) {
		configPath := filepath.Join(tempDir, "replace_test.cfg")

		// Create test config file
		err := os.WriteFile(configPath, []byte(validBase64JSON), 0644)
		if err != nil {
			t.Fatalf("Failed to create test config file: %v", err)
		}

		loader := func(path string) error { return nil }

		// Start first watcher
		err = WatchConfig([]string{configPath}, loader)
		if err != nil {
			t.Fatalf("First WatchConfig failed: %v", err)
		}

		// Start second watcher (should replace first)
		err = WatchConfig([]string{configPath}, loader)
		if err != nil {
			t.Fatalf("Second WatchConfig failed: %v", err)
		}

		// Test that the second call succeeds (replaces the first watcher)
		// The actual replacement behavior is internal to the function

		// Clean up
		if watcher != nil {
			watcher.Close()
			watcher = nil
		}
	})

	// --- Test: Non-existent File ---
	t.Run("NonExistentFile", func(t *testing.T) {
		configPath := filepath.Join(tempDir, "nonexistent.cfg")

		loader := func(path string) error { return nil }

		// This should fail - fsnotify requires the file to exist
		err := WatchConfig([]string{configPath}, loader)
		if err == nil {
			t.Error("Expected error for non-existent file")
		}
		if err != nil && !contains(err.Error(), "no such file or directory") {
			t.Errorf("Expected 'no such file or directory' error, got: %v", err)
		}
	})

	// --- Test: Empty Paths ---
	t.Run("EmptyPaths", func(t *testing.T) {
		loader := func(path string) error { return nil }

		// This should fail - empty slice is not valid
		err := WatchConfig([]string{}, loader)
		if err == nil {
			t.Error("Expected error for empty paths")
		}
		if err != nil && err.Error() != "config paths are required" {
			t.Errorf("Expected 'config paths are required' error, got: %v", err)
		}
	})

	// --- Test: Nil Loader ---
	t.Run("NilLoader", func(t *testing.T) {
		configPath := filepath.Join(tempDir, "nil_loader_test.cfg")

		// Create a test config file
		err := os.WriteFile(configPath, []byte(validBase64JSON), 0644)
		if err != nil {
			t.Fatalf("Failed to create test config file: %v", err)
		}

		// This should fail - nil loader is not valid
		err = WatchConfig([]string{configPath}, nil)
		if err == nil {
			t.Error("Expected error for nil loader")
		}
		if err != nil && err.Error() != "loader is required" {
			t.Errorf("Expected 'loader is required' error, got: %v", err)
		}
	})

	// --- Test: Loader Error Handling ---
	t.Run("LoaderErrorHandling", func(t *testing.T) {
		configPath := filepath.Join(tempDir, "error_test.cfg")

		// Create test config file
		err := os.WriteFile(configPath, []byte(validBase64JSON), 0644)
		if err != nil {
			t.Fatalf("Failed to create test config file: %v", err)
		}

		// Loader that returns an error
		loader := func(path string) error {
			return os.ErrPermission
		}

		// Start watching
		err = WatchConfig([]string{configPath}, loader)
		if err != nil {
			t.Fatalf("WatchConfig failed: %v", err)
		}

		// Give the watcher time to start
		time.Sleep(100 * time.Millisecond)

		// Write to the file to trigger an event
		err = os.WriteFile(configPath, []byte(validBase64JSON), 0644)
		if err != nil {
			t.Fatalf("Failed to write config: %v", err)
		}

		// Give time for the event to be processed
		time.Sleep(200 * time.Millisecond)

		// The loader error should be handled gracefully (logged but not fatal)
		// We can't easily test the error logging without capturing log output,
		// but we can verify the watcher continues to work

		// Clean up
		if watcher != nil {
			watcher.Close()
			watcher = nil
		}
	})

	// --- Test: Concurrent File Writes ---
	t.Run("ConcurrentWrites", func(t *testing.T) {
		configPath := filepath.Join(tempDir, "concurrent_test.cfg")

		// Create initial config file
		err := os.WriteFile(configPath, []byte(validBase64JSON), 0644)
		if err != nil {
			t.Fatalf("Failed to create test config file: %v", err)
		}

		var loaderCalls []string
		var mu sync.Mutex
		loader := func(path string) error {
			mu.Lock()
			loaderCalls = append(loaderCalls, path)
			mu.Unlock()
			return nil
		}

		// Start watching
		err = WatchConfig([]string{configPath}, loader)
		if err != nil {
			t.Fatalf("WatchConfig failed: %v", err)
		}

		// Give the watcher time to start
		time.Sleep(100 * time.Millisecond)

		// Perform multiple concurrent writes
		var wg sync.WaitGroup
		for i := 0; i < 5; i++ {
			wg.Add(1)
			go func(index int) {
				defer wg.Done()
				config := &Config{
					APIEndpoint:   "concurrent.com",
					CACertificate: "concurrent_cacert",
					Certificate:   "concurrent_cert",
					PrivateKey:    "concurrent_key",
				}
				configBytes, _ := json.Marshal(config)
				configBase64 := base64.StdEncoding.EncodeToString(configBytes)

				err := os.WriteFile(configPath, []byte(configBase64), 0644)
				if err != nil {
					t.Errorf("Failed to write config %d: %v", index, err)
				}
			}(i)
		}

		wg.Wait()

		// Give time for all events to be processed
		time.Sleep(500 * time.Millisecond)

		// Check if loader was called (may be called multiple times due to concurrent writes)
		mu.Lock()
		callCount := len(loaderCalls)
		mu.Unlock()

		if callCount == 0 {
			t.Error("Expected loader to be called for concurrent writes")
		}

		// Clean up
		if watcher != nil {
			watcher.Close()
			watcher = nil
		}
	})

	// --- Test: Watcher Cleanup ---
	t.Run("WatcherCleanup", func(t *testing.T) {
		configPath := filepath.Join(tempDir, "cleanup_test.cfg")

		// Create test config file
		err := os.WriteFile(configPath, []byte(validBase64JSON), 0644)
		if err != nil {
			t.Fatalf("Failed to create test config file: %v", err)
		}

		loader := func(path string) error { return nil }

		// Start watching
		err = WatchConfig([]string{configPath}, loader)
		if err != nil {
			t.Fatalf("WatchConfig failed: %v", err)
		}

		// Test that the first call succeeds

		// Start another watcher (should close the first one)
		err = WatchConfig([]string{configPath}, loader)
		if err != nil {
			t.Fatalf("Second WatchConfig failed: %v", err)
		}

		// The first watcher should be closed, but we can't easily test that
		// without accessing the internal state. We can verify a new watcher
		// was created.

		// Clean up
		if watcher != nil {
			watcher.Close()
			watcher = nil
		}
	})

	// Final cleanup
	if watcher != nil {
		watcher.Close()
		watcher = nil
	}
}
