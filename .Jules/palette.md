## 2024-03-24 - [Micro-UX Improvements to TUI Config]
**Learning:** TUI configuration interfaces often lack visual affordances found in GUI settings menus. Users need immediate feedback on what actions are available and the boundaries of valid inputs. Adding explicit bounds checking (clamping) and visual indicators for boolean states (ON/OFF vs True/False) significantly improves confidence.
**Action:** When designing TUI configuration screens, always include:
1. Visible key hints (e.g., "[←/→] Adjust").
2. Explicit visual states for toggles (Color-coded ON/OFF).
3. "Invisible" safety rails (clamping values) to prevent invalid configurations.

## 2024-05-22 - [Contextual Help in TUI Lists]
**Learning:** In compact TUI menus, complex settings often lack space for inline explanations. A "selected item description" pattern—rendering helper text for the active item in a dedicated footer area—provides essential context without cluttering the UI. This significantly aids accessibility for new users.
**Action:** When creating TUI selection lists for configuration or complex actions, implement a parallel `Descriptions` slice and render the active description dynamically below the list.
