## 2024-05-22 - Sentinel initialized

## 2026-01-29 - Hardcoded RPC Credentials
**Vulnerability:** Found hardcoded API keys for Shyft and Helius RPC services directly embedded in `config.yaml` URLs.
**Learning:** Hardcoded secrets often hide in "convenience" configurations like full URLs. Developers might paste the full connection string to get it working quickly.
**Prevention:** Enforce environment variable injection for all credential-bearing parameters. Use helper methods to construct sensitive URLs at runtime rather than storing them whole.
