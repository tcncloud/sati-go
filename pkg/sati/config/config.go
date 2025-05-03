// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//
// Copyright 2024 TCN Inc

package config // <--- Changed from 'sati'

import (
	"encoding/base64"
	"encoding/json"
	"os"
)

type Config struct {
	CACertificate           string `json:"ca_certificate"`
	Certificate             string `json:"certificate"`
	PrivateKey              string `json:"private_key"`
	FingerprintSHA256       string `json:"fingerprint_sha256"`
	FingerprintSHA256String string `json:"fingerprint_sha256_string"`
	APIEndpoint             string `json:"api_endpoint"`
	CertificateName         string `json:"certificate_name"`
	CertificateDescription  string `json:"certificate_description"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	decoded := make([]byte, base64.StdEncoding.DecodedLen(len(data)))
	n, err := base64.StdEncoding.Decode(decoded, data)
	if err != nil {
		return nil, err
	}
	// DEBUG: Print decoded bytes before unmarshal
	// fmt.Printf("LoadConfig decoded bytes: %s\n", string(decoded[:n]))
	var config Config
	if err := json.Unmarshal(decoded[:n], &config); err != nil {
		return nil, err
	}
	return &config, nil
}

// NewConfigFromString creates a Config object from a base64 encoded JSON string.
func NewConfigFromString(configString string) (*Config, error) {
	decoded, err := base64.StdEncoding.DecodeString(configString)
	if err != nil {
		return nil, err
	}
	// DEBUG: Print decoded bytes before unmarshal
	// fmt.Printf("NewConfigFromString decoded bytes: %s\n", string(decoded))
	var config Config
	if err := json.Unmarshal(decoded, &config); err != nil {
		return nil, err
	}
	return &config, nil
}
