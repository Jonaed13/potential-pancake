## 2024-05-22 - Sentinel initialized

## 2024-05-23 - Hardcoded Secrets in Config
**Vulnerability:** Found hardcoded API keys for Shyft and Helius directly in `config/config.yaml` URLs.
**Learning:** Keys were embedded in URL strings to simplify passing them to `RPCClient`, but this bypassed environment variable separation.
**Prevention:** Always use config helpers to inject secrets from environment variables into URLs at runtime, never commit them to YAML/JSON.
