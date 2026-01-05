## 2024-05-23 - Hardcoded API Keys Removed
**Vulnerability:** Hardcoded Helius and Shyft API keys were found in `config/config.yaml`.
**Learning:** Developers often commit configuration files with secrets for convenience during initial setup, forgetting that these files are part of the repository.
**Prevention:**
1.  Always use `.env.example` for templates.
2.  Use environment variables for ALL secrets.
3.  Inject secrets into configuration at runtime (e.g., constructing URLs with keys) rather than storing full URLs with keys in config files.
4.  Add pre-commit hooks to scan for high-entropy strings or known key patterns.
