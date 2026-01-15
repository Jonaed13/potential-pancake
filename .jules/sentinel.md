## 2024-03-24 - API Key Injection in RPC URLs
**Vulnerability:** Hardcoded API keys were found in `config/config.yaml` for both Helius and Shyft RPC endpoints.
**Learning:** Config files (even if not committed initially) can easily be committed by mistake or leak secrets if used as templates. Relying on a single `apiKey` field in the RPC client struct is insufficient when using multiple providers with different auth mechanisms (header vs query param).
**Prevention:**
1.  Always use environment variables for secrets.
2.  If providers require different auth methods (e.g., query param for Helius, header for Shyft), construct the full authenticated URL at runtime in the main entry point, rather than baking incomplete auth logic into a generic client.
3.  Use `.env.example` to document required keys without exposing actual values.
