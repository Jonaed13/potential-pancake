## 2026-01-29 - [Batch RPC Calls and Program IDs]
**Learning:** `getTokenAccountsByOwner` RPC method filters by `programId`. Batch fetching all token accounts for an owner using the standard Token Program ID will NOT return accounts for other programs (like Token-2022).
**Action:** When using batch RPC calls for optimization, always implement a fallback mechanism (individual fetch) for items missing from the batch result, as they might belong to a different program or standard not covered by the batch query.
