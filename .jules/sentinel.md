
## 2024-05-22 - Sentinel initialized

## 2024-05-22 - Hardcoded API Keys in Config
**Vulnerability:** Hardcoded Helius and Shyft API keys were found in `config/config.yaml` URLs.
**Learning:** Example configuration files often inadvertently become production configuration files, leading to secret leakage. Even if intended for "ease of use", secrets should never be committed.
**Prevention:** Use `.env` files for all secrets and reference them via environment variables in the code. Ensure `.gitignore` excludes `.env` and verify no secrets are in committed `*.yaml` or `*.json` files.
