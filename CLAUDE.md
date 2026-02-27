# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build & Development Commands

| Command | Description |
|---|---|
| `make build` | Build binary to `/tmp/bin/GoYoutubeDL` |
| `make run` | Build and run the server (listens on `:9001`) |
| `make run/live` | Hot-reload dev server via Air |
| `make test` | Run all tests: `go test -v -race -buildvcs ./...` |
| `make test/cover` | Run tests with coverage and open HTML report |
| `make lint` | Format code and run golangci-lint |
| `make audit` | Full quality check: vet, staticcheck, govulncheck, race tests |
| `make tidy` | Format code and tidy go.mod |

Run a single test: `go test -v -run TestName ./internal/client/...`

## External Requirement

**ffmpeg** must be installed on the host. It is used for merging audio+video streams (via kkdai/youtube `DownloadComposite`) and embedding metadata into `.m4a` audio files.

## Architecture

Three-layer dependency injection pattern with interfaces at each boundary:

```
main.go (wiring + Echo server setup)
  └── handler.Handler  →  depends on handler.Servicer interface
        └── service.Service  →  depends on service.YtClient interface
              └── client.Client  →  wraps client.ytdlr interface (kkdai/youtube)
```

Each layer defines its own interface for the layer below it (`internal/{layer}/interface.go`), keeping dependencies decoupled and independently testable.

### Package Responsibilities

- **`internal/client`** — Wraps the `kkdai/youtube` library. Handles video/playlist fetching, file streaming, thumbnail extraction, and audio metadata embedding via ffmpeg.
- **`internal/service`** — Business logic: URL validation, playlist detection, ffmpeg availability check, orchestrates client calls.
- **`internal/handler`** — Echo HTTP handlers. Renders HTMX partial responses using Go templates. Serves the download UI and media player.
- **`internal/models`** — Shared types (`Video`, `Playlist`, `Image`, `Player`), custom error types, template renderer, and constants (e.g., `DirPath` for the Downloads directory).

### Frontend

HTMX-driven UI with Go HTML templates in `html/`. Static assets (htmx.min.js, favicons) in `assets/`. Templates use `{{define}}` blocks to serve partial HTML responses for HTMX swaps.

### Key Routes

- `GET /` — Main download page
- `POST /getInfo` — Fetch video/playlist metadata from YouTube URL
- `GET /info` — Get quality options for a video ID
- `POST /download` — Download video or audio
- `GET /player` — Media player listing downloaded files
- `GET /resource/*` — Serves downloaded media files from `./Downloads/`

## Testing Conventions

- Table-driven tests throughout
- Mocks generated with `gomock` from interface files (`internal/client/mock_interface.go` is generated from `internal/client/interface.go`)
- Assertions use `github.com/stretchr/testify/assert`
- Tests exist for client and service layers; handler layer has no tests yet

## Linting

Strict golangci-lint v2 configuration (`.golangci.yml`): max line length 140, cyclomatic/cognitive complexity limit 10, function length limit 100 lines/65 statements. Several linters are relaxed for test files.
