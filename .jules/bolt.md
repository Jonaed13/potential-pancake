## 2025-01-26 - Solana RPC Batching
**Learning:** Fetching token accounts on Solana requires querying both the Legacy Token Program and Token-2022 Program to ensure complete state. Batching these into concurrent calls (or a single filtered call if supported) is critical for performance when tracking multiple positions.
**Action:** Always check both Program IDs when fetching "all" token accounts. Use `GetAllTokenAccounts` pattern instead of iterating per mint.
