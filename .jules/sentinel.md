## 2024-05-22 - Sentinel initialized

## 2024-05-23 - Hardcoded Secrets in Config URLs
**Vulnerability:** Found Helius and Shyft API keys hardcoded directly into the URL strings in `config/config.yaml`.
**Learning:** Even when `env` variable support exists (like `shyft_api_key_env`), developers might still hardcode keys in URLs if the URL format requires the key as a query parameter and no helper exists to inject it.
**Prevention:** Do not allow URLs in config to contain query parameters with secrets. Implement strict helper methods (like `GetShyftRPCURL`) that take a base URL and inject the secret from the environment at runtime.
