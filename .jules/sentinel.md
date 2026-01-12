## 2024-03-24 - [CRITICAL] Hardcoded API Keys in Config and Tools
**Vulnerability:** Found hardcoded Helius and Shyft API keys in `config/config.yaml` and several helper tools in `tools/`.
**Learning:** Developers often commit secrets in "test" or "helper" scripts, thinking they won't be used in production. Config files with "default" values are also a common trap.
**Prevention:**
1. Use `.env` files for ALL secrets, even in development.
2. Use `git-secrets` or pre-commit hooks to scan for high-entropy strings.
3. Never hardcode keys in `config.yaml` defaults; use environment variable injection (like `SHYFT_API_KEY_ENV`).
