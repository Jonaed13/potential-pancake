## 2025-02-18 - Removal of Hardcoded Secrets in Tools
**Vulnerability:** Found multiple standalone tools (`tools/`) containing hardcoded private keys (Solana wallets) and API keys (Helius, Shyft).
**Learning:** Developers often treat "test" or "benchmark" tools as safe zones for hardcoded secrets, assuming they won't be compiled into the main binary. However, committing these files to the repo exposes the secrets to anyone with read access.
**Prevention:** Strictly prohibit hardcoded secrets in ALL files, including tools and tests. Use environment variables or command-line flags for secret injection. Enforce pre-commit hooks to scan for high-entropy strings or known key patterns.
