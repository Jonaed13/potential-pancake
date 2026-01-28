## 2024-03-24 - [Micro-UX Improvements to TUI Config]
**Learning:** TUI configuration interfaces often lack visual affordances found in GUI settings menus. Users need immediate feedback on what actions are available and the boundaries of valid inputs. Adding explicit bounds checking (clamping) and visual indicators for boolean states (ON/OFF vs True/False) significantly improves confidence.
**Action:** When designing TUI configuration screens, always include:
1. Visible key hints (e.g., "[←/→] Adjust").
2. Explicit visual states for toggles (Color-coded ON/OFF).
3. "Invisible" safety rails (clamping values) to prevent invalid configurations.

## 2025-05-23 - [Static Indicators in "Live" Dashboards]
**Learning:** Hardcoding status indicators (like "[LIVE]") in dashboard headers creates a dangerous disconnect when the underlying system state (Paused/Running) changes. Users trust the header as the source of truth for system status.
**Action:** Always bind status headers directly to the state variable (`m.Running`), using both text changes ("LIVE" vs "PAUSED") and color changes (Green vs Red) to ensure the status is unambiguous.
