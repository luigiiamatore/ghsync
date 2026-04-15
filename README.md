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
│   └── report_clear.go   # clear reports command
├── internal/
│   └── report/
│       └── report.go     # sync report model and storage
├── main.go
└── README.md
```

## Development

### Running Tests

Tests are included to verify core functionality across commands:

- **Authentication tests**: verify token is saved with correct permissions
- **Pagination tests**: verify repository sync with 500+ repositories

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

## License

MIT
