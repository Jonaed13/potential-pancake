## 2024-05-22 - Sentinel initialized
## 2025-05-23 - URL Secret Injection Pattern
**Vulnerability:** Hardcoded API keys in `config.yaml` URLs were used because the `RPCClient` assumed a single authentication method (header-based) which didn't work for all providers (e.g., Helius requires query param).
**Learning:** When integrating multiple third-party providers with different auth mechanisms (header vs query param), the configuration system must support dynamic injection of secrets into URLs at runtime rather than relying on static strings.
**Prevention:** Use `fmt.Sprintf` or `net/url` to construct URLs with secrets injected from environment variables at the point of client initialization, ensuring secrets are never committed to config files.
