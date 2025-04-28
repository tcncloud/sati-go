# Specification

- Single binary application
- Go project structure follows best practices:
    - CLI entry point and cobra command implementation in `cmd/sati-client/main.go`
    - All business logic, config loading, and client setup in `pkg/sati/`
    - Command implementations in `pkg/cmd/` (one file per command)
    - Generated proto code in `internal/genproto/`
    - Proto definitions in `tcnapi/`
- Uses spf13/cobra for CLI
- Reads a config file (TOML or base64-encoded JSON) with the following fields:
    - ca_certificate / ca_file
    - certificate / cert_file
    - private_key / key_file
    - fingerprint_sha256
    - fingerprint_sha256_string
    - api_endpoint
    - certificate_name
    - certificate_description
- Connects to `api_endpoint` using `certificate`, `private_key`, and `ca_certificate` for mTLS
- Exposes cobra subcommands for GateService gRPC methods, including:
    - `dial`: Initiate a call
    - `put-call-on-simple-hold`: Put a call on hold
    - `take-call-off-simple-hold`: Take a call off hold
    - `stop-call-recording`: Stop call recording
    - `get-recording-status`: Get call recording status
    - `rotate-certificate`: Rotate certificate
    - `log`: Send a log message
    - `update-scrub-list-entry`: Update a scrub list entry
    - `add-scrub-list-entries`: Add entries to a scrub list
    - `remove-scrub-list-entries`: Remove entries from a scrub list
    - `submit-job-results`: Submit job results (with support for oneof result types)
    - (Extensible: more GateService methods can be added as subcommands)
- Uses generated gRPC client code from local proto definitions
- All dependencies and proto imports are resolved locally or via Go modules
- See [README.md](README.md) for usage, configuration, and command examples.


