# ghsync TODO

## Alta Priorità

- [ ] **Proper error handling** - Log errors to stderr instead of stdout
  - Use `fmt.Fprintf(os.Stderr, ...)` for errors
  - Keep stdout clean for normal output

- [ ] **Extract constants** - Remove hardcoded `.ghsync` paths
  - Define `const ConfigDir = ".ghsync"`
  - Use throughout codebase

- [ ] **Quiet git output** - Suppress verbose git clone/pull output
  - Add `--quiet` flag to git commands
  - Or redirect stdout/stderr

## Media Priorità

- [ ] **Configuration file support** - Move from just token to full config
  - Create `.ghsync/config.yml` instead of just `config`
  - Store: token, default dir, exclude patterns, etc.
  - Parse YAML on startup

- [ ] **Add --quiet flag** - For script usage
  - Disable box output in quiet mode
  - Only print essentials

- [ ] **Timeout git operations** - Prevent hanging on large repos
  - Add context timeout to git clone/pull commands
  - Configurable timeout via config file

## Nice to Have

- [ ] **Parallel cloning/pulling** - Use goroutines for speed
  - Implement worker pool pattern
  - Limit concurrent operations (e.g., max 5 at a time)

- [ ] **Structured logging** - Better log organization
  - Consider minimal logging setup
  - Log levels (info, warn, error)

- [ ] **GitHub Actions CI/CD** - Automated testing on push
  - Run tests on PR
  - Build and test on multiple OS

- [ ] **--dry-run flag** - Preview what would happen
  - Show which repos would be cloned vs pulled
  - No actual git operations

- [ ] **Filtering support** - Include/exclude repo patterns
  - `--include "pattern"` for specific repos
  - `--exclude "pattern"` to skip repos

## Completed ✓

- [x] Pagination support for 100+ repositories
- [x] Progress bar during sync
- [x] Beautiful box-formatted output
- [x] Report reading and displaying
- [x] Report clearing with confirmation
- [x] Token authentication
- [x] Test suite (pull, auth, report, clear)
- [x] Error collection and display
