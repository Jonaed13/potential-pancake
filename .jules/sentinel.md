## 2024-05-22 - Sentinel initialized

## 2024-05-22 - Config Secrets Refactor
**Vulnerability:** Hardcoded API keys found in `config/config.yaml` URL parameters.
**Learning:** Config files often accumulate secrets over time. Centralizing URL construction logic (e.g., `GetShyftRPCURL`) in the config package prevents duplication and ensures consistent secret injection.
**Prevention:** Use a helper function to construct URLs with secrets, never embed them in static config files.
