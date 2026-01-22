# Economist TUI

Terminal UI and CLI to browse and read The Economist.

## Install

### Homebrew (macOS)

```bash
brew install --HEAD https://raw.githubusercontent.com/tmustier/economist-tui/main/Formula/economist-tui.rb
```

### From source

```bash
git clone https://github.com/tmustier/economist-tui
cd economist-tui
make install
```

**Prereqs:** Go 1.25+, Chrome/Chromium (for login and article fetching).

## Quick start

```bash
# Interactive browsing (TUI)
economist browse finance

# Login (one-time, for full articles)
economist login

# Read an article
economist read "https://www.economist.com/finance-and-economics/2026/01/19/article-slug"

# Headlines for scripts
economist headlines leaders --json
```

## Commands

- `browse [section]` — interactive TUI
  - `↑/↓` navigate, `←/→` page, `Enter` read, `b` back, `c` columns, `Esc` clear, `q` quit
- `headlines [section]` — list headlines
  - `-n/--number`, `-s/--search`, `--json`, `--plain`
- `read [url|-]` — read full article (`--raw`, `--wrap`, `--columns`)
- `login` — open browser to authenticate
- `serve` — warm browser daemon (`--status`, `--stop`)
- `sections` — list sections

Global flags: `--version`, `--debug`, `--no-color`

## Configuration

Config + cookies: `~/.config/economist-tui/` (migrates from `economist-cli`)
Cache: `~/.config/economist-tui/cache` (1h TTL)

## Notes

- RSS provides ~300 items per section (~10 months)
- Full articles require an active Economist subscription

## License

MIT
