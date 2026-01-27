## 2024-05-22 - Sentinel initialized

## 2024-05-22 - Missing Security Tests
**Vulnerability:** Critical security controls (Rate Limiting) and their verification tests (`server_security_test.go`) were documented in memory but missing from the codebase.
**Learning:** Do not trust memory or documentation blindly; verify the existence of security controls in the actual code.
**Prevention:** Always perform a gap analysis between documented security features and implemented code.
