## 2026-01-19 - RPC Parsing Performance
**Learning:** `fmt.Sscanf` is significantly slower (~43x) and allocates more memory than `strconv.ParseUint` for simple integer parsing. In high-frequency RPC response parsing loops (like `GetTokenAccountsByOwner`), this overhead adds up.
**Action:** Prefer `strconv` functions over `fmt` scanning for simple data extraction in critical paths.
