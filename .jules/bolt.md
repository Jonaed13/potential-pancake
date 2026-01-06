# Bolt's Journal ⚡

## 2024-05-22 - TUI Rendering Optimization
**Learning:** In Bubble Tea apps, heavy logic inside `Update()` or `View()` blocks the UI thread. Moving calculation to separate goroutines and using `tea.Cmd` is crucial for responsiveness.
**Action:** Always check if heavy operations in TUI models are blocking the render loop.

## 2024-05-22 - JSON Decoding
**Learning:** `json.NewDecoder` is generally faster and uses less memory than `json.Unmarshal` for large or streaming data, but for small payloads, the difference is negligible and `Unmarshal` is often more readable. However, strictly using `NewDecoder` for http responses is a good habit.
**Action:** Prefer streaming decoders for API responses.

## 2024-05-23 - RPC Client Pooling
**Learning:** Creating a new `http.Client` for every request kills performance due to lack of connection pooling. Reusing a global or shared client is essential.
**Action:** Ensure HTTP clients are long-lived.
