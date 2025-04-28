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

package sati

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/url"
	"strings"

	"io"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func SetupClient(cfg *Config) (*grpc.ClientConn, error) {
	cert, err := tls.X509KeyPair([]byte(cfg.Certificate), []byte(cfg.PrivateKey))
	if err != nil {
		return nil, fmt.Errorf("failed to load client cert: %w", err)
	}
	caCertPool := x509.NewCertPool()
	if ok := caCertPool.AppendCertsFromPEM([]byte(cfg.CACertificate)); !ok {
		return nil, fmt.Errorf("failed to append CA cert")
	}
	creds := credentials.NewTLS(&tls.Config{
		ServerName:   "exile-proxy", // TODO: remove this once we have a proper TLS certificate
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
	})

	endpoint := ParseAPIEndpoint(cfg.APIEndpoint)

	conn, err := grpc.Dial(endpoint, grpc.WithTransportCredentials(creds))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to API: %w", err)
	}
	return conn, nil
}

func ParseAPIEndpoint(raw string) string {
	if len(raw) == 0 {
		return raw
	}
	if u, err := url.Parse(raw); err == nil && u.Host != "" {
		host := u.Host
		if u.Scheme == "https" && !strings.Contains(host, ":") {
			host += ":443"
		}
		return host
	}
	if strings.HasPrefix(raw, "http://") {
		host := strings.TrimPrefix(raw, "http://")
		return host
	}
	if strings.HasPrefix(raw, "https://") {
		host := strings.TrimPrefix(raw, "https://")
		if !strings.Contains(host, ":") {
			host += ":443"
		}
		return host
	}
	return raw
}

// IsStreamEnd returns true if the error indicates the end of a gRPC stream.
func IsStreamEnd(err error) bool {
	return err == io.EOF
}
