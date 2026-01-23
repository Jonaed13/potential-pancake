## 2026-01-23 - Batch RPC Calls for Token Balances
**Learning:** The application was making N separate RPC calls to `getTokenAccountsByOwner` (filtered by mint) for every active position in the `monitorPositions` loop. This caused linear latency scaling with the number of positions.
**Action:** Use `getTokenAccountsByOwner` with `programId` filter to fetch ALL token accounts in a single (or two, for Token-2022) RPC call, then map them in memory. This reduces complexity from O(N) to O(1) network calls per tick.
