# ghsync

CLI tool to sync your GitHub repositories and generate a report after the synchronization.

## Installation

```bash
git clone https://github.com/luigiiamatore/ghsync.git
cd ghsync
go build -o ghsync
```

## Usage

Authenticate with your GitHub token:

```bash
./ghsync auth
```

Sync all your repositories:

```bash
./ghsync pull
./ghsync pull --dir /path/to/repos  # custom directory
```

View the latest sync report:

```bash
./ghsync report
./ghsync report --all     # view all reports
```

Clear all sync reports:

```bash
./ghsync report clear     # removes all stored reports
```

## Project Structure

```text
ghsync/
├── cmd/
│   ├── auth.go           # authentication command
│   ├── auth_test.go      # authentication tests
│   ├── pull.go           # sync command
│   ├── pull_test.go      # sync command tests
│   ├── report.go         # report command
│   ├── report_test.go    # report command tests
│   ├── report_clear.go   # clear reports command
│   └── report_clear_test.go # clear reports command tests
├── internal/
│   └── report/
│       └── report.go     # sync report model and storage
├── main.go
└── README.md
```

## Development

### Running Tests

Tests are included to verify core functionality across all commands:

- **Authentication tests**: verify token is saved with correct permissions
- **Pagination tests**: verify repository sync with 500+ repositories
- **Report tests**: verify report files are correctly read and parsed
- **Clear reports tests**: verify safe deletion of report directories

Run all tests:

```bash
go test ./...
```

Run tests with coverage report:

```bash
go test -cover ./...
```

Run pagination test (verifies sync with 500+ repositories):

```bash
go test -run TestPullPaginationWithMockServer ./cmd -v
```

Run authentication test (verifies token save with permissions):

```bash
go test -run TestSaveToken ./cmd -v
```

Run report tests (verifies report reading and parsing):

```bash
go test -run TestReadSyncReport ./cmd -v
```

Run clear reports test (verifies safe deletion):

```bash
go test -run TestClear ./cmd -v
```

## License

MIT
