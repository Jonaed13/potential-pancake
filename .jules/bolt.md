## 2025-05-18 - RPC Response Parsing
**Learning:** `fmt.Sscanf` is significantly slower (~36x) than `strconv.ParseUint` for parsing simple integer strings from JSON-RPC responses. In high-frequency trading bots, this parsing overhead adds up when processing large numbers of token accounts or balance updates.
**Action:** Always prefer `strconv` functions over `fmt` scanning for critical path parsing, especially in RPC client implementations.
