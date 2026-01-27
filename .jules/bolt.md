## 2024-05-23 - [Tools Maintenance]
**Learning:** Changes to internal packages (like `token` or `blockchain`) can break CLI tools in `tools/` which import them. These tools are often overlooked during standard `go test ./internal/...` runs.
**Action:** Always run `go build ./...` to catch breaking changes in `tools/` when modifying shared internal packages.
