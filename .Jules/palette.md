## 2024-03-24 - [Micro-UX Improvements to TUI Config]
**Learning:** TUI configuration interfaces often lack visual affordances found in GUI settings menus. Users need immediate feedback on what actions are available and the boundaries of valid inputs. Adding explicit bounds checking (clamping) and visual indicators for boolean states (ON/OFF vs True/False) significantly improves confidence.
**Action:** When designing TUI configuration screens, always include:
1. Visible key hints (e.g., "[←/→] Adjust").
2. Explicit visual states for toggles (Color-coded ON/OFF).
3. "Invisible" safety rails (clamping values) to prevent invalid configurations.

## 2025-05-23 - Contextual Help in Config Modals
**Learning:** TUI configuration settings (like "MinEntryPercent") can be ambiguous. Users benefit significantly from immediate, in-context explanations of what each setting controls, rather than relying on external documentation.
**Action:** When creating configuration interfaces, always pair the setting name with a brief, faint/italicized description that explains its impact or units.
