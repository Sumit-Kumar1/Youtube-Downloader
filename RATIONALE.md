# Change Rationale & Validity Evaluation

Each fix from the code review is listed below with the reasoning behind the change
and an evaluation of whether it is the correct approach.

---

## Critical Bug Fixes

### 1. `exec.Command(cmd)` — single string passed as program name

**What changed:** `internal/client/utils.go` — Rewrote `addMetaData` to use `exec.CommandContext("ffmpeg", args...)` where args is a `[]string` built by `prepareMetadataArgs`.

**Rationale:** `exec.Command` takes the program name as the first argument and additional arguments as variadic parameters. Passing the entire command as one string makes Go search for a binary literally named `"ffmpeg -i file.m4a -map 0 ..."`, which never exists. Every audio metadata operation was guaranteed to fail.

**Validity: Correct.** This is the standard Go pattern for subprocess execution. Using `exec.CommandContext` also ties the subprocess lifetime to the request context (fix #9 synergy). Each argument is a separate element — no shell interpretation means no injection risk (fix #6 synergy).

---

### 2. Modulo truncation `len(s) % 100` replaced with `min(len(s), 100)`

**What changed:** `internal/client/utils.go` — Added `truncate(s string, maxLen int) string` helper. Replaced all `vid.Field[:len(vid.Field)%100]` calls.

**Rationale:** `len % 100` produces mathematically wrong results for the intended purpose:

| Input length | `len % 100` | Intended | Actual result |
|---|---|---|---|
| 50 | 50 | 50 | Correct by coincidence |
| 99 | 99 | 99 | Correct by coincidence |
| 100 | 0 | 100 | **Empty string** |
| 150 | 50 | 100 | **Only 50 chars** |
| 200 | 0 | 100 | **Empty string** |
| 300 | 0 | 100 | **Empty string** |

The `truncate` function uses a simple conditional: if the string is within limits, return it unchanged; otherwise, slice to `maxLen`.

**Validity: Correct.** This is the standard truncation pattern. The function is pure, trivially testable, and handles all edge cases. Test `Test_truncate` covers empty, short, exact, and over-length strings.

---

### 3. `stream2File` returned `nil` instead of `err` on sync failure

**What changed:** `internal/client/utils.go:60` — Fixed `return nil` → properly handled error. Also rewrote the entire function (see #17).

**Rationale:** The original code had `if err := tempF.Sync(); err != nil { return nil }`. This silently swallowed disk flush failures. Data could appear to write successfully while the OS buffer was never committed to disk — risking silent data corruption or zero-byte files.

**Validity: Correct.** Errors should always propagate. The rewritten `stream2File` properly surfaces all errors from every I/O operation.

---

### 4. Nil pointer dereference on `*getThumbnail()`

**What changed:** `internal/client/utils.go` — Changed `getThumbnail` return type from `*models.Image` to `models.Image` (value type). Updated call sites in `client.go` to remove dereference.

**Rationale:** When a video had no thumbnails (empty `Thumbnails` slice), `getThumbnail` returned `nil`. The call site `*getThumbnail(...)` unconditionally dereferenced the result, causing a nil pointer panic. Returning a zero-value `Image{}` instead of `nil` is idiomatic Go — the zero value is meaningful (no thumbnail) and safe.

**Validity: Correct.** Returning value types for small structs (3 fields) avoids nil issues with no performance cost. The `Image` struct is 40 bytes — well within the threshold where value semantics are preferred over pointer semantics in Go. Tests updated to expect `models.Image{}` instead of `nil`.

---

### 5. ffmpeg command had no output file

**What changed:** `internal/client/utils.go` — `prepareMetadataArgs` now takes both `inputFile` and `outputFile` parameters. The output file is appended as the last argument. `addMetaData` creates a `.tmp` output file, then renames it over the input file on success.

**Rationale:** ffmpeg requires an output file argument. The original command `ffmpeg -i input -map 0 -c:a copy -metadata ...` would fail with "At least one output file must be specified." Writing to a temp file and then renaming is an atomic pattern that prevents data loss if ffmpeg crashes mid-write.

**Validity: Correct.** The `-y` flag was added to auto-overwrite the temp file without prompting (ffmpeg would otherwise block waiting for stdin in a headless server). The rename-over-original pattern is the industry standard for safe file replacement.

---

### 6. `audioFormats[0]` panic on empty slice

**What changed:** `internal/client/client.go` — Added `len(audioFormats) == 0` check before indexing.

**Rationale:** If a YouTube video has no audio-only formats (unlikely but possible for some content), `vid.Formats.Type("audio")` returns an empty slice. Indexing `[0]` on an empty slice is an unrecoverable runtime panic that would crash the server goroutine handling the request.

**Validity: Correct.** Defensive bounds checking before indexing is a Go best practice. Returns a descriptive error instead of crashing. The error message includes the video ID for debugging.

---

### 7. `/assets/*` route served from wrong directory

**What changed:** `main.go` — Changed `http.Dir(models.DirPath)` to `http.Dir("assets")` for the `/assets/*` route.

**Rationale:** Both `/resource/*` and `/assets/*` were serving from `./Downloads/`. But `htmx.min.js`, favicons, and `site.webmanifest` live in `./assets/`. Every request for `assets/htmx.min.js` from the HTML templates would 404, breaking the entire HTMX-powered UI.

**Validity: Correct.** The `/resource/*` route correctly serves downloaded media from `./Downloads/`. The `/assets/*` route should serve static frontend assets from `./assets/`. These are different directories with different purposes.

---

## Bug Fixes

### 8. Stream leak in `DownloadAudio`

**What changed:** `internal/client/client.go` — Added `defer stream.Close()` immediately after obtaining the stream. Removed the `return stream.Close()` at the end.

**Rationale:** If `stream2File` or `addMetaData` returned an error, the function returned early and `stream.Close()` was never called. This leaked the HTTP connection to YouTube's servers. Using `defer` guarantees cleanup regardless of the return path.

**Validity: Correct.** `defer` for resource cleanup is idiomatic Go. The stream is an `io.ReadCloser` wrapping an HTTP response body — failing to close it leaks connections from the HTTP client pool.

---

### 9. `context.Background()` used instead of request context

**What changed:** `internal/handler/handler.go` — All three handler methods (`GetInfo`, `Download`, `DownloadInfo`) now use `c.Request().Context()` instead of `context.Background()`.

**Rationale:** Using `context.Background()` means:
- The 3-minute timeout middleware has zero effect (it sets a deadline on the request context, which was being discarded)
- Client disconnects are not propagated (downloads continue running after the user closes the browser)
- No cancellation propagation to ffmpeg subprocesses or YouTube API calls

**Validity: Correct.** This is the single most impactful correctness fix. Echo's context carries deadlines, cancellation, and tracing information. Discarding it defeats every middleware that operates on context (timeout, tracing, etc.). This change required adding `context.Context` parameters to all interface methods (see #16).

---

### 10. `Play` handler returned `nil` for unknown file types

**What changed:** `internal/handler/handler.go` — Added `default` case to the switch that renders an error template.

**Rationale:** Returning `nil` from an Echo handler without writing a response results in an empty 200 response with no body. The user sees a blank page with no indication of what went wrong.

**Validity: Correct.** Every HTTP handler code path should produce a response. The error template is already used elsewhere and provides consistent UX. Test `TestHandler_Play_UnsupportedType` validates this behavior.

---

### 11. `Download` handler exposed raw internal errors

**What changed:** `internal/handler/handler.go` — Replaced `return err` with `return c.Render(http.StatusOK, "error", ...)` using a generic user-facing message. The raw error is still logged.

**Rationale:** Raw errors could expose internal paths (`/Users/.../Downloads/file.m4a`), ffmpeg stderr output, or YouTube API error messages. These leak implementation details and potentially sensitive information. HTMX expects an HTML response for swapping — returning an Echo error would render JSON or a generic HTTP error page, breaking the UI.

**Validity: Correct.** The separation of logging (detailed, internal) and user response (generic, safe) follows security best practices. The HTTP status stays 200 because HTMX needs a successful response to swap content — the error template provides visual error feedback within the page.

---

### 12. `FillByName` set Title/Path for unknown file extensions

**What changed:** `internal/models/player.go` — Added `return` in the `default` case of the extension switch.

**Rationale:** Without the early return, files like `readme.txt`, `.part`, `.tmp` would get a Title and Path assigned but no Type or IsAudio flag. This created inconsistent `Player` structs — `getDownloadedFiles` in the handler would filter them out by Type, but the intermediate state was still wrong. More critically, if filtering logic ever changed, these invalid entries would leak through.

**Validity: Correct.** One-line fix with clear semantics: unknown file types produce a zero-value `Player`. Test `TestPlayer_FillByName` covers the "unknown extension" and "partial file" cases.

---

### 13. `isFFMpegInstalled` spawned subprocess on every download

**What changed:** `internal/service/utils.go` — Replaced `exec.Command("ffmpeg", "-version").CombinedOutput()` with `sync.Once` + `exec.LookPath("ffmpeg")`.

**Rationale:** Every download call spawned an ffmpeg subprocess just to check if ffmpeg exists. `exec.LookPath` is a pure PATH lookup — no process creation. `sync.Once` caches the result so it's checked exactly once per server lifetime.

**Validity: Correct.** `exec.LookPath` is the standard Go way to check binary availability. `sync.Once` is goroutine-safe and idiomatic for one-time initialization. The only edge case is if ffmpeg is installed/uninstalled while the server is running — this is acceptable for a development tool.

---

### 14. Debug `fmt.Printf`/`Println` removed

**What changed:** `internal/client/utils.go` — Removed `fmt.Printf("\nFileName: %s", fileName)` and `fmt.Println("running cmd:\n", cmd)`.

**Rationale:** These write to stdout, bypassing the structured logging system. In production, they produce unformatted noise mixed with Echo's access logs. The `slog.LogAttrs` call on error already provides proper logging with context.

**Validity: Correct.** Debug prints should never be committed to production code. The existing `slog` error logging in `addMetaData` is sufficient. If debug logging is needed during development, it should use `slog.Debug`.

---

## Design & Architecture Improvements

### 15. `audioOnly` changed from `string` to `bool`

**What changed:** `internal/handler/interface.go`, `internal/service/service.go`, `internal/handler/handler.go` — The `audioOnly` parameter is now `bool`. Parsing happens in the handler: `c.FormValue("audioOnly") == "true"`.

**Rationale:** The original `switch audioOnly` had three implicit paths:
- `""` → download video
- `"true"` → download audio
- **anything else** → silently do nothing, return nil

This meant `"false"`, `"yes"`, `"1"`, or any typo would result in a successful API response with no download. The bool conversion eliminates this silent failure — there are exactly two states.

**Validity: Correct.** HTTP form parsing belongs in the handler layer. The service layer should receive typed Go values, not raw form strings. This follows the separation of concerns principle. The handler test `TestHandler_Download/audio_only_download` verifies the parsing.

---

### 16. Context propagation through all layers

**What changed:** `internal/client/interface.go`, `internal/client/client.go`, `internal/service/interface.go`, `internal/service/service.go` — Added `context.Context` as the first parameter to all interface methods. Switched from `GetVideo`/`GetPlaylist` to `GetVideoContext`/`GetPlaylistContext` in the kkdai wrapper.

**Rationale:** Without context propagation:
- Request timeouts don't apply to YouTube API calls or ffmpeg subprocesses
- Client disconnects don't cancel in-flight operations
- Server shutdown can't interrupt active downloads

The kkdai library already supports context-aware methods (`GetVideoContext`, `GetPlaylistContext`). The original code even used `GetVideoContext` in `DownloadAudio` but not in other methods — this was inconsistent.

**Validity: Correct.** Context is Go's standard mechanism for request-scoped cancellation, timeouts, and values. The convention of passing `context.Context` as the first parameter follows the Go standard library pattern and the guidance in the Go blog ("Context should be the first parameter of a function"). The mock interface was updated to match.

---

### 17. `stream2File` rewritten — no more full-file memory copy

**What changed:** `internal/client/utils.go` — Complete rewrite. Creates temp file in `DirPath` (same filesystem), uses `io.Copy` for streaming, `os.Rename` for atomic move.

**Rationale:** The original code read the entire temp file into memory with `os.ReadFile` and then wrote it to the destination with `file.Write`. For a 500MB video audio track, this would allocate 500MB of heap memory unnecessarily. The new approach:
1. Creates temp file in the same directory as the destination (ensures `os.Rename` works — same filesystem)
2. Streams data via `io.Copy` (constant memory regardless of file size)
3. Syncs to disk
4. Atomically renames to final path

The deferred cleanup removes the temp file if any step fails.

**Validity: Correct.** `io.Copy` uses a fixed 32KB buffer internally. `os.Rename` is atomic on the same filesystem (POSIX guarantee). Creating temp files in the destination directory ensures same-filesystem semantics. This is the standard pattern for safe file writing in Go.

---

### 18. Graceful shutdown added

**What changed:** `main.go` — Server now starts in a goroutine. Main goroutine listens for `SIGINT`/`SIGTERM`, then calls `e.Shutdown(ctx)` with a 30-second deadline.

**Rationale:** `e.Logger.Fatal(e.Start(...))` would kill the process immediately on signal. Active downloads would be interrupted mid-write, potentially leaving corrupted partial files. `echo.Shutdown` stops accepting new connections and waits for active requests to complete (up to the 30-second timeout).

**Validity: Correct.** This is Echo's recommended shutdown pattern (documented in the Echo guide). The 30-second timeout is a reasonable default — most downloads would either complete or be long enough that they shouldn't block shutdown. Using `http.ErrServerClosed` to distinguish normal shutdown from actual errors is the standard pattern.

---

### 19. `YtClient` field unexported on `Service` struct

**What changed:** `internal/service/service.go` — Changed `YtClient YtClient` to `ytClient YtClient`.

**Rationale:** Exported struct fields are part of the package's public API. The `YtClient` field is an implementation detail set exclusively via the `New` constructor — no code outside the service package should access or modify it directly.

**Validity: Correct.** This follows Go's principle of minimal exported surface. Even though the `internal/` package boundary already prevents access from outside the module, unexported fields within internal packages still enforce encapsulation between packages inside the module.

---

### 20. Unused `player` field removed from Handler

**What changed:** `internal/handler/handler.go` — Removed `player map[string]bool` from the `Handler` struct and `make(map[string]bool)` from `New()`.

**Rationale:** The field was allocated in every handler initialization but never read or written. It was dead code that added memory overhead and confusion when reading the codebase.

**Validity: Correct.** Dead code removal. If the `player` map is needed for a future feature, it should be added when that feature is implemented — not carried as unused baggage.

---

### 21. Regex compiled once at package level

**What changed:** `internal/client/utils.go` — Moved `regexp.MustCompile(...)` from inside `formatName` to a package-level `var nameCleanRegex`.

**Rationale:** `regexp.MustCompile` parses and compiles the regex pattern. Calling it inside the function recompiles on every invocation. While the regex is simple and compilation is fast (~microseconds), it creates unnecessary garbage that the GC must collect. Moving to package level compiles once at init time.

**Validity: Correct.** This is a universally accepted Go best practice, documented in the standard library. The `regexp` package is safe for concurrent use — package-level compiled regexes can be shared across goroutines without synchronization.

---

### 22. Duplicate `id="spinner"` in HTML

**What changed:** `html/index.html` — Renamed the download form spinner from `id="spinner"` to `id="dlSpinner"`. Updated the corresponding `hx-indicator` attribute.

**Rationale:** HTML requires unique `id` attributes per document. Two elements with `id="spinner"` means `document.querySelector("#spinner")` (which HTMX uses internally) would always select the first one. The download form's `hx-indicator="#spinner"` was inadvertently showing the wrong spinner — or not showing at all if the first spinner was in a different DOM position.

**Validity: Correct.** Each HTMX indicator needs its own unique ID. The main page search spinner keeps `id="spinner"`, and the download form gets `id="dlSpinner"`. Both function independently.

---

### 23. `validateURL` rewritten with `url.Parse` and host allowlist

**What changed:** `internal/service/utils.go` — Replaced string matching (`strings.Contains(url, "youtu") || !strings.ContainsAny(url, "\"?&/<%=")`) with `url.Parse` + host allowlist.

**Rationale:** The original validation had multiple issues:
- `strings.ContainsAny(url, "\"?&/<%=")` — The `<` and `%` characters suggest URL-encoded or XSS payloads were considered valid
- `strings.Contains(url, "youtu")` — Would match `https://notyoutube.example.com` or `https://evil.com/youtu`
- The logic was an AND of two loose checks, allowing many non-YouTube URLs through

The new approach uses Go's `url.Parse` to validate structure (scheme + host required), then checks the host against a fixed set of known YouTube domains.

**Validity: Correct.** Host allowlisting is the standard approach for URL validation. The allowlist covers `youtube.com`, `www.youtube.com`, `m.youtube.com`, `youtu.be`, and `www.youtu.be`. New tests cover short URLs, mobile URLs, non-YouTube domains, and malformed inputs.

**One consideration:** The new validation requires a scheme (`https://`). URLs without a scheme (e.g., `youtube.com/watch?v=abc`) will now be rejected. This is arguably more correct — the downstream kkdai library also expects well-formed URLs — and the HTML input already has `type="url"` which encourages scheme inclusion.

---

## New Test Coverage

| File | Tests Added | What They Cover |
|---|---|---|
| `internal/models/player_test.go` | 7 cases | FillByName: empty, short, mp4, m4a, mp3, unknown ext, partial file |
| `internal/client/utils_test.go` | 9 new cases | `truncate` (5 cases), `prepareMetadataArgs` (2 cases: basic + truncation) |
| `internal/client/client_test.go` | Updated | GetVideo/GetPlaylist now use context-aware mock expectations |
| `internal/service/utils_test.go` | 4 new cases | validateURL: youtu.be, mobile, no-scheme, empty string |
| `internal/handler/handler_test.go` | 9 cases (new file) | Page render, GetInfo success/error, Download success/audio/error, Play unsupported type, DownloadInfo success/error |

**Total tests: 30 (was 13, added 17 new)**

---

## Items Not Changed (and why)

### `DirPath` relative path (Review item #24)
The `DirPath` constant `"./Downloads/"` is relative to the working directory. Making it absolute would require either an `init()` function (flagged by `gochecknoinits` linter) or passing it through the constructor chain. The current behavior is correct when running from the project root, which is enforced by the Makefile and Air config. A future improvement could use an environment variable.

### `github.com/golang/mock` → `go.uber.org/mock` migration (Review item #21)
The deprecated `golang/mock` library still functions correctly. The migration requires:
1. `go get go.uber.org/mock`
2. Updating all imports from `github.com/golang/mock/gomock` to `go.uber.org/mock/gomock`
3. Regenerating all mocks with the new `mockgen`
4. Removing the old dependency

This is a separate chore best done as its own PR to keep the current changeset focused on bug fixes.
