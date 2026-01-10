## 2024-03-24 - [Micro-UX Improvements to TUI Config]
**Learning:** TUI configuration interfaces often lack visual affordances found in GUI settings menus. Users need immediate feedback on what actions are available and the boundaries of valid inputs. Adding explicit bounds checking (clamping) and visual indicators for boolean states (ON/OFF vs True/False) significantly improves confidence.
**Action:** When designing TUI configuration screens, always include:
1. Visible key hints (e.g., "[←/→] Adjust").
2. Explicit visual states for toggles (Color-coded ON/OFF).
3. "Invisible" safety rails (clamping values) to prevent invalid configurations.

## 2024-05-24 - [Contextual Help in Modals]
**Learning:** In keyboard-driven interfaces, users often hesitate to change settings because they don't want to lose their place or guess what a cryptic acronym means. Immediate, selection-aware help text reduces this friction significantly without cluttering the main view.
**Action:** For lists of settings or complex options, always reserve a footer area for "Contextual Help" that updates instantly as the user navigates through the items.
