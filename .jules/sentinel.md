## 2024-05-22 - Sentinel initialized

## 2024-05-22 - [Secrets in Config]
**Vulnerability:** Hardcoded API keys found in `config.yaml` URLs (Shyft and Helius).
**Learning:** Even "public" keys in URLs are dangerous if committed. Configuration files often inadvertently host secrets when URLs are constructed statically.
**Prevention:** Use environment variables for all secrets. Construct URLs dynamically in code (`main.go`) by injecting keys from env vars, rather than storing full URLs with secrets in config files. Added `FallbackAPIKeyEnv` to support multiple providers (Shyft/Helius) securely.
