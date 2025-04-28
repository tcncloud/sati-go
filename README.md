# Sati Client

A command-line client for interacting with the TCN Exile Gate gRPC API. This tool allows you to manage agents, calls, scrub lists, job results, and more from the terminal.

## Features
- Agent management (list, update, upsert, etc.)
- Call control (dial, hold, record, etc.)
- Scrub list management
- Job result submission
- Organization and configuration queries

## Prerequisites
- Go 1.20 or newer
- Access to a running TCN Exile Gate gRPC server
- Proper configuration file (see below)

## Installation
Clone the repository and build the binary:

```sh
go install github.com/tcncloud/sati-go/cmd/sati-client@latest
```
or 
```
git clone <repo-url>
cd sati-go
go build -o sati-client ./cmd/sati-client
```

## Configuration
Generate and copy the configuration certificate from operator.tcn.com and save it into a local file.


## Usage
Run the client with a command and specify the config file:

```sh
./sati-client <command> --config com.tcn.exiles.sati.config.cfg [flags]
```

## Available Commands (examples)
- `dial` — Initiate a call:
  ```sh
  ./sati-client dial --partner-agent-id AGENT_ID --phone-number 1234567890 --config com.tcn.exiles.sati.config.cfg
  ```
- `put-call-on-simple-hold` — Put a call on hold:
  ```sh
  ./sati-client put-call-on-simple-hold --partner-agent-id AGENT_ID --config com.tcn.exiles.sati.config.cfg
  ```
- `add-scrub-list-entries` — Add entries to a scrub list:
  ```sh
  ./sati-client add-scrub-list-entries --scrub-list-id LIST_ID --entries '[{"content":"1234567890","notes":"spam"}]' --config com.tcn.exiles.sati.config.cfg
  ```
- `remove-scrub-list-entries` — Remove entries from a scrub list:
  ```sh
  ./sati-client remove-scrub-list-entries --scrub-list-id LIST_ID --entries "1234567890,0987654321" --config com.tcn.exiles.sati.config.cfg
  ```
- `submit-job-results` — Submit job results:
  ```sh
  ./sati-client submit-job-results --job-id JOB_ID --result '{"error_result":{"message":"fail"}}' --config com.tcn.exiles.sati.config.cfg
  ```

## Help
For a full list of commands and flags, run:

```sh
./sati-client --help
./sati-client <command> --help
```

## License
Apache 2.0. See [LICENSE](LICENSE) for details. 