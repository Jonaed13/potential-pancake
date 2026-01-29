## 2024-03-24 - [Micro-UX Improvements to TUI Config]
**Learning:** TUI configuration interfaces often lack visual affordances found in GUI settings menus. Users need immediate feedback on what actions are available and the boundaries of valid inputs. Adding explicit bounds checking (clamping) and visual indicators for boolean states (ON/OFF vs True/False) significantly improves confidence.
**Action:** When designing TUI configuration screens, always include:
1. Visible key hints (e.g., "[←/→] Adjust").
2. Explicit visual states for toggles (Color-coded ON/OFF).
3. "Invisible" safety rails (clamping values) to prevent invalid configurations.

## 2024-05-22 - [Contextual Help in Configuration]
**Learning:** In complex TUI applications, users often struggle to remember what specific acronyms or settings do without leaving the interface. Providing inline descriptions for the currently selected item transforms a "guessing game" into an informed decision-making process.
**Action:** Implement dynamic description footers in lists/modals that update instantly as the user navigates, providing "just-in-time" documentation.
