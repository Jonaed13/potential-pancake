## 2025-05-23 - [Optimization: fmt.Sscanf vs strconv.ParseUint]
**Learning:** In hot RPC response parsing paths (like `GetTokenAccountsByOwner`), using `fmt.Sscanf` was found to be ~25x slower than `strconv.ParseUint` (1148ns vs 45ns). While convenient, `fmt` reflection overhead is significant in loops.
**Action:** For high-frequency parsing of simple types (int/uint), always prefer `strconv`. When retrofitting, ensure error handling behavior (like defaulting to 0) is preserved to avoid regressions.
