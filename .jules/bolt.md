## 2026-01-18 - fmt.Sscanf vs strconv.ParseUint
**Learning:** `fmt.Sscanf` is ~26x slower (730ns vs 27ns) than `strconv.ParseUint` for simple integer parsing.
**Action:** Use `strconv` for parsing API responses in high-frequency loops (like RPC updates).
