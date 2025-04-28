# Specification

- Single binary application
- Go project structure follows best practices:
    - CLI entry point and cobra command implementation in `cmd/sati-client/main.go`
    - All business logic, config loading, and client setup in `pkg/sati/`
    - Generated proto code in `internal/genproto/`
    - Proto definitions in `tcnapi/`
- Uses spf13/cobra for CLI
- Reads a base64-encoded JSON config file with the following fields:
    - ca_certificate
    - certificate
    - private_key
    - fingerprint_sha256
    - fingerprint_sha256_string
    - api_endpoint
    - certificate_name
    - certificate_description
- Connects to `api_endpoint` using `certificate`, `private_key`, and `ca_certificate` for mTLS
- Exposes cobra subcommands for GateService gRPC methods:
    - `get-client-config`: Calls GateService.GetClientConfiguration and prints the result
    - `get-org-info`: Calls GateService.GetOrganizationInfo and prints the result
    - (Extensible: more GateService methods can be added as subcommands)
- Uses generated gRPC client code from local proto definitions
- All dependencies and proto imports are resolved locally or via Go modules


