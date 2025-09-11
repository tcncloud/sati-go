package config

import (
	"encoding/base64"
	"os"
	"path/filepath"
	"testing"
)

const validBase64JSON = "eyJhcGlfZW5kcG9pbnQiOiJ0ZXN0LmNvbSIsImNhX2NlcnRpZmljYXRlIjoidGVzdF9jYSIsImNlcnRpZmljYXRlIjoidGVzdF9jZXJ0IiwicHJpdmF0ZV9rZXkiOiJ0ZXN0X2tleSJ9Cg=="

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

		// Note: The current implementation has a defer watcher.Close() which closes
		// the watcher immediately when WatchConfig returns, so we can't test the
		// watcher state directly. We test the function's return value instead.
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

		// Simple loader function
		loader := func(path string) error {
			return nil
		}

		// Start watching multiple paths
		err = WatchConfig([]string{configPath1, configPath2}, loader)
		if err != nil {
			t.Fatalf("WatchConfig failed: %v", err)
		}

		// Note: The current implementation has a defer watcher.Close() which closes
		// the watcher immediately when WatchConfig returns, so we can't test the
		// watcher state directly. We test the function's return value instead.
	})

	// --- Test: Empty Paths ---
	t.Run("EmptyPaths", func(t *testing.T) {
		// Simple loader function
		loader := func(path string) error {
			return nil
		}

		// Start watching with empty paths
		err := WatchConfig([]string{}, loader)
		if err != nil {
			t.Fatalf("WatchConfig failed: %v", err)
		}

		// Note: The current implementation has a defer watcher.Close() which closes
		// the watcher immediately when WatchConfig returns, so we can't test the
		// watcher state directly. We test the function's return value instead.
	})

	// --- Test: Non-existent File ---
	t.Run("NonExistentFile", func(t *testing.T) {
		configPath := filepath.Join(tempDir, "nonexistent.cfg")

		// Simple loader function
		loader := func(path string) error {
			return nil
		}

		// Start watching non-existent file
		err := WatchConfig([]string{configPath}, loader)
		if err != nil {
			t.Fatalf("WatchConfig failed: %v", err)
		}

		// Note: The current implementation has a defer watcher.Close() which closes
		// the watcher immediately when WatchConfig returns, so we can't test the
		// watcher state directly. We test the function's return value instead.
	})
}
