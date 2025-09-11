package config

import (
	"encoding/base64"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

// --- Test Data ---
// Use a struct to define the expected config for easier comparison.
var expectedValidConfig = &Config{
	APIEndpoint:   "test.com",
	CACertificate: "cacert",
	Certificate:   "cert",
	PrivateKey:    "key",
	// Fingerprint/Name/Desc fields omitted as they are not in the minimal valid JSON
}

// Marshal the expected config to get the canonical JSON string.
var validJSONBytes, _ = json.Marshal(expectedValidConfig)

// var validJSONString = string(validJSONBytes) // Unused variable

// Base64 encode the canonical JSON string.
var validBase64JSON = base64.StdEncoding.EncodeToString(validJSONBytes)

// Other test cases.
const (
	invalidJSON         = `{"api_endpoint":"test.com",}` // Malformed JSON
	invalidBase64       = "%%%invalid&&&"
	partiallyValidJSON  = `{"api_endpoint":"partial.com"}` // Valid JSON, but missing fields
	base64PartiallyJSON = "eyJhcF9lbmRwb2ludCI6InBhcnRpYWwuY29tIn0="
)

// --- Helper Functions ---.
func checkConfigFields(t *testing.T, cfg *Config, expected *Config) {
	t.Helper()

	if cfg == nil {
		t.Fatal("Config object is nil")
	}
	// Using reflect.DeepEqual for a comprehensive check
	if !reflect.DeepEqual(cfg, expected) {
		t.Errorf("Config mismatch:\n Got: %+v\nWant: %+v", cfg, expected)
	}
	// Individual checks for easier debugging if DeepEqual fails
	if cfg.APIEndpoint != expected.APIEndpoint {
		t.Errorf("Expected APIEndpoint '%s', got '%s'", expected.APIEndpoint, cfg.APIEndpoint)
	}

	if cfg.CACertificate != expected.CACertificate {
		t.Errorf("Expected CACertificate '%s', got '%s'", expected.CACertificate, cfg.CACertificate)
	}

	if cfg.Certificate != expected.Certificate {
		t.Errorf("Expected Certificate '%s', got '%s'", expected.Certificate, cfg.Certificate)
	}

	if cfg.PrivateKey != expected.PrivateKey {
		t.Errorf("Expected PrivateKey '%s', got '%s'", expected.PrivateKey, cfg.PrivateKey)
	}
}

// --- Tests ---

func TestLoadConfig(t *testing.T) {
	tempDir := t.TempDir()

	// --- Subtest: Valid Config File ---
	t.Run("ValidFile", func(t *testing.T) {
		validFilePath := filepath.Join(tempDir, "valid.cfg")

		err := os.WriteFile(validFilePath, []byte(validBase64JSON), 0644)
		if err != nil {
			t.Fatalf("Failed to write valid test config file: %v", err)
		}

		cfg, err := LoadConfig(validFilePath)
		if err != nil {
			t.Fatalf("LoadConfig failed for valid file: %v", err)
		}

		checkConfigFields(t, cfg, expectedValidConfig)
	})

	// --- Subtest: Non-existent File ---
	t.Run("NonExistentFile", func(t *testing.T) {
		_, err := LoadConfig(filepath.Join(tempDir, "nonexistent.cfg"))
		if err == nil {
			t.Error("LoadConfig succeeded for non-existent file, expected error")
		}
	})

	// --- Subtest: Invalid Base64 File ---
	t.Run("InvalidBase64", func(t *testing.T) {
		invalidBase64FilePath := filepath.Join(tempDir, "invalid_base64.cfg")

		err := os.WriteFile(invalidBase64FilePath, []byte(invalidBase64), 0644)
		if err != nil {
			t.Fatalf("Failed to write invalid base64 test file: %v", err)
		}

		_, err = LoadConfig(invalidBase64FilePath)
		if err == nil {
			t.Error("LoadConfig succeeded for invalid base64 file, expected error")
		}
	})

	// --- Subtest: Invalid JSON File (after base64 decode) ---
	t.Run("InvalidJSON", func(t *testing.T) {
		invalidJSONFilePath := filepath.Join(tempDir, "invalid_json.cfg")
		// Explicitly encode the invalid JSON to be sure
		encodedInvalidJSON := base64.StdEncoding.EncodeToString([]byte(invalidJSON))

		err := os.WriteFile(invalidJSONFilePath, []byte(encodedInvalidJSON), 0644)
		if err != nil {
			t.Fatalf("Failed to write invalid json test file: %v", err)
		}

		_, err = LoadConfig(invalidJSONFilePath)
		if err == nil {
			t.Error("LoadConfig succeeded for invalid json file, expected error")
		}
	})
}

func TestNewConfigFromString(t *testing.T) {
	// --- Subtest: Valid String ---
	t.Run("ValidString", func(t *testing.T) {
		cfg, err := NewConfigFromString(validBase64JSON)
		if err != nil {
			t.Fatalf("NewConfigFromString failed for valid string: %v", err)
		}

		checkConfigFields(t, cfg, expectedValidConfig)
	})

	// --- Subtest: Invalid Base64 String ---
	t.Run("InvalidBase64String", func(t *testing.T) {
		_, err := NewConfigFromString(invalidBase64)
		if err == nil {
			t.Error("NewConfigFromString succeeded for invalid base64 string, expected error")
		}
	})

	// --- Subtest: Invalid JSON String (after base64 decode) ---
	t.Run("InvalidJSONString", func(t *testing.T) {
		// Explicitly encode the invalid JSON to be sure
		encodedInvalidJSON := base64.StdEncoding.EncodeToString([]byte(invalidJSON))

		_, err := NewConfigFromString(encodedInvalidJSON)
		if err == nil {
			t.Error("NewConfigFromString succeeded for invalid json string, expected error")
		}
	})
}
