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
```

## Project Structure

```text
ghsync/
├── cmd/
│   ├── auth.go       # authentication command
│   ├── pull.go       # sync command
│   └── report.go     # report command
├── internal/
│   └── report/
│       └── report.go # sync report model and storage
├── main.go
└── README.md
```

## License

MIT
