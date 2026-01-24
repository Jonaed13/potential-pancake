## 2025-05-20 - [fmt.Sscanf vs strconv.ParseUint]
**Learning:** `fmt.Sscanf` uses reflection and is significantly slower (~27x) than `strconv.ParseUint` for simple integer parsing. In hot loops like `GetTokenAccountsByOwner`, this adds up.
**Action:** Prefer `strconv` for parsing primitive types in performance-critical paths.
