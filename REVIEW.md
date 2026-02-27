# Code Review — Youtube-Downloader

Reviewed on the `audio-metadata` branch. Findings are grouped by severity.

---

## Critical Bugs (will crash or silently fail at runtime)

### 1. `exec.Command` receives entire command as program name
**File:** `internal/client/utils.go:100`

`exec.Command(cmd)` passes the full command string as a single argument (the program name). `exec.Command` expects the program as the first arg and the remaining arguments separately. This will always fail with "executable file not found".

**Fix:** Split into program + args:
```go
func addMetaData(c context.Context, vid *youtube.Video, fileName string) error {
    fileName = filepath.Clean(models.DirPath + "/" + fileName)

    args := prepareMetadataArgs(vid, fileName)

    out, err := exec.CommandContext(c, "ffmpeg", args...).CombinedOutput()
    if err != nil {
        slog.LogAttrs(c, slog.LevelError, "ffmpeg metadata failed",
            slog.String("output", string(out)),
            slog.String("error", err.Error()))
        return err
    }
    return nil
}

func prepareMetadataArgs(vid *youtube.Video, inputFile string) []string {
    outputFile := inputFile + ".tmp"

    args := []string{"-i", inputFile, "-map", "0", "-c:a", "copy"}

    meta := map[string]string{
        "title":       truncate(vid.Title, 100),
        "author":      truncate(vid.Author, 100),
        "description": truncate(vid.Description, 100),
        "year":        strconv.Itoa(vid.PublishDate.Year()),
    }

    for k, v := range meta {
        args = append(args, "-metadata", k+"="+v)
    }

    args = append(args, outputFile)
    return args
}
```

### 2. Modulo truncation produces wrong lengths
**File:** `internal/client/utils.go:113-115`

```go
vid.Title[:len(vid.Title)%100]
```

This does **not** truncate to 100 chars. The modulo produces incorrect results:
- 50-char string → `50 % 100 = 50` (ok by accident)
- 200-char string → `200 % 100 = 0` → **empty string**
- 300-char string → `300 % 100 = 0` → **empty string**
- 150-char string → `150 % 100 = 50` → takes only 50 chars

**Fix:** Use `min()`:
```go
func truncate(s string, maxLen int) string {
    if len(s) <= maxLen {
        return s
    }
    return s[:maxLen]
}
```

### 3. `stream2File` swallows sync error — returns `nil` instead of `err`
**File:** `internal/client/utils.go:60`

```go
if err := tempF.Sync(); err != nil {
    return nil  // BUG: should be return err
}
```

Data can silently fail to flush to disk.

### 4. Nil pointer dereference when video has no thumbnails
**File:** `internal/client/client.go:35, 68`

`getThumbnail()` returns `nil` when the thumbnail slice is empty. The call site dereferences it:
```go
Thumbnail: *getThumbnail(ytVid.Thumbnails),  // panics if nil
```

**Fix:** Return a zero-value `Image` instead of nil from `getThumbnail`, or check before dereferencing.

### 5. `ffmpeg` command has no output file
**File:** `internal/client/utils.go:110`

The command is built as `ffmpeg -i <input> -map 0 -c:a copy -metadata ...` but never specifies an output file. `ffmpeg` requires an output file argument at the end. The command will fail even after fixing the `exec.Command` issue.

### 6. `audioFormats[0]` panics on empty slice
**File:** `internal/client/client.go:119`

If the video has no audio-only formats, `vid.Formats.Type("audio")` returns an empty slice and `audioFormats[0]` panics.

**Fix:** Check length before indexing:
```go
audioFormats := vid.Formats.Type("audio")
if len(audioFormats) == 0 {
    return fmt.Errorf("no audio formats available for video %s", id)
}
```

### 7. `/assets/*` route serves from wrong directory
**File:** `main.go:49-50`

```go
e.GET("/assets/*", echo.WrapHandler(http.StripPrefix("/assets/",
    http.FileServer(http.Dir(models.DirPath)))))
```

This serves `/assets/*` from `./Downloads/`, but `htmx.min.js`, favicons, and `site.webmanifest` live in `./assets/`. The HTML templates reference `assets/htmx.min.js`, etc. — these requests will 404.

**Fix:**
```go
e.GET("/assets/*", echo.WrapHandler(http.StripPrefix("/assets/",
    http.FileServer(http.Dir("assets")))))
```

---

## Bugs (incorrect behavior, not necessarily crash)

### 8. Stream leak on error in `DownloadAudio`
**File:** `internal/client/client.go:105-133`

If `stream2File` or `addMetaData` returns an error, the function returns early and `stream.Close()` on line 132 is never called. Use `defer` immediately after obtaining the stream.

**Fix:**
```go
stream, _, err := c.ytd.GetStreamContext(ctx, vid, &audioFormats[0])
if err != nil {
    return err
}
defer stream.Close()
```
Then remove the final `return stream.Close()` and just `return nil` (or return the deferred close error using a named return).

### 9. `context.Background()` used instead of request context
**File:** `internal/handler/handler.go:63, 79, 91`

All three handler methods create `context.Background()` instead of using the request's context. This means:
- The 3-minute timeout middleware has no effect on downstream operations
- Request cancellations (client disconnects) are ignored
- Downloads continue running even after the client disconnects

**Fix:** Use `c.Request().Context()` everywhere.

### 10. `handler.Play` returns `nil` for unknown file types
**File:** `internal/handler/handler.go:42-57`

If `p.Type` is neither `AudioType` nor `VideoType`, the handler returns `nil` without writing any HTTP response. This results in an empty 200 response.

**Fix:** Return an error or a "not supported" render after the switch.

### 11. `handler.Download` exposes raw errors to clients
**File:** `internal/handler/handler.go:83`

`return err` passes internal errors directly to Echo, which renders them as 500 responses with the raw error message. Internal errors (file paths, ffmpeg output, etc.) should not leak to users.

**Fix:** Log the error, return a user-friendly rendered error response (same pattern as `GetInfo`).

### 12. `FillByName` sets Title/Path for unknown extensions
**File:** `internal/models/player.go:23-41`

When the extension doesn't match (not m4a/mp3/mp4), the `switch` falls through to default (which does nothing), but lines 39-40 still execute unconditionally, setting `Title` and `Path`. This means non-media files (e.g., `.txt`, `.part`) get a Title and Path but no Type, creating an inconsistent state.

**Fix:** Move lines 39-40 inside the recognized cases, or add a `return` in the default case.

### 13. `isFFMpegInstalled` spawns a subprocess on every download
**File:** `internal/service/utils.go:27-33`

Every call to `Download` runs `ffmpeg -version` as a subprocess. This is wasteful — check once at startup and cache the result.

### 14. Debug print statements left in production code
**File:** `internal/client/utils.go:89, 98`

```go
fmt.Printf("\nFileName: %s", fileName)
fmt.Println("running cmd:\n", cmd)
```

These should be removed or replaced with structured logging (`slog`).

---

## Design & Architecture Improvements

### 15. `audioOnly` parameter is a string, should be bool
**File:** `internal/service/service.go:33`

`Download(ctx context.Context, id, qual, audioOnly string)` — the `audioOnly` parameter is compared against `"true"` and `""`. Parse it to a `bool` in the handler layer before passing it to the service.

Any value other than `""` and `"true"` (e.g., `"yes"`, `"1"`) silently results in no download at all — the function returns `nil` without doing anything.

### 16. No context propagation through client layer
**Files:** `internal/client/client.go`, `internal/service/interface.go`

`GetVideo`, `GetPlaylist`, and `GetDownloadInfo` don't accept `context.Context`. This prevents cancellation, tracing, and timeout propagation. The underlying kkdai library supports context — `GetVideoContext` is already used in `DownloadAudio` but not in other methods.

### 17. `stream2File` is inefficient — double memory copy
**File:** `internal/client/utils.go:47-84`

The function:
1. Streams data to a temp file via `io.Copy`
2. Syncs the temp file
3. Reads the **entire** temp file into memory with `os.ReadFile`
4. Writes it to the final destination

Step 3 loads the entire file into RAM. For a large audio file, this is wasteful. Use `os.Rename` if on the same filesystem, or use `io.Copy` from the temp file to the destination.

### 18. No graceful shutdown
**File:** `main.go:52`

`e.Logger.Fatal(e.Start(":9001"))` doesn't handle OS signals. Active downloads will be interrupted on shutdown. Use `e.Shutdown(ctx)` with signal handling.

### 19. Exported field `YtClient` on `Service` struct
**File:** `internal/service/service.go:10`

```go
type Service struct {
    YtClient YtClient
}
```

This is an internal package, but the field should still be unexported (`ytClient`) since it's an implementation detail that should only be set via the constructor.

### 20. `handler.player` field is unused
**File:** `internal/handler/handler.go:14`

The `player map[string]bool` field is initialized in `New()` but never read or written anywhere. Remove it.

### 21. Deprecated mock library
**File:** `go.mod:7`

`github.com/golang/mock` is archived/deprecated. Migrate to `go.uber.org/mock` which is the maintained successor.

### 22. `formatName` recompiles regex on every call
**File:** `internal/client/utils.go:41`

`regexp.MustCompile` is called inside the function body, recompiling the regex every time. Move it to a package-level `var`.

### 23. Duplicate HTML `id="spinner"`
**File:** `html/index.html:62, 131`

Two elements share `id="spinner"` — the main page loading spinner and the download form spinner. `hx-indicator="#spinner"` on line 41 targets the first one, but the form on line 117 also references `#spinner`. Duplicate IDs cause unpredictable behavior. Use distinct IDs.

### 24. `DirPath` is a relative path
**File:** `internal/models/player.go:4`

`"./Downloads/"` depends on the working directory at runtime. If the binary is launched from a different directory, downloads go to the wrong place. Consider resolving to an absolute path at startup.

### 25. `validateURL` is fragile
**File:** `internal/service/utils.go:16`

```go
if !strings.Contains(url, "youtu") || !strings.ContainsAny(url, "\"?&/<%=") {
```

This allows through any URL containing "youtu" and any of those special characters. The `<` and `%` checks are concerning — they hint at encoded/malicious input. Consider using `url.Parse` and checking the host against known YouTube domains (`youtube.com`, `youtu.be`, `m.youtube.com`).

---

## Summary

| Severity | Count |
|----------|-------|
| Critical (crash/silent failure) | 7 |
| Bug (incorrect behavior) | 7 |
| Design improvement | 11 |

The audio-metadata feature (`addMetaData` / `prepareMetadataCommand`) has the most critical issues — the `exec.Command` usage, modulo truncation, missing output file, and lack of shell-safe argument handling mean it will not work at all in its current form. Fixing those is the immediate priority.
