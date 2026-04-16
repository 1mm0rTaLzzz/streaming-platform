# Design System Strategy: The Stadium Atmosphere

## 1. Overview & Creative North Star
**Creative North Star: "The Digital Arena"**

This design system is not a static webpage; it is a high-stakes, broadcast-quality experience. We are moving away from the "standard website" feel to embrace **Kinetic Editorialism**. The goal is to capture the raw energy of the World Cup 2026—the scale of the USA, the passion of Mexico, and the precision of Canada.

We achieve this through **intentional asymmetry** and **tonal depth**. By utilizing oversized, aggressive typography juxtaposed with refined, glass-like data modules, we create a "Premium Broadcast" aesthetic. The layout should feel like a live scoreboard—dynamic, urgent, and authoritative. We break the grid using overlapping elements (e.g., a player’s silhouette breaking out of a `surface-container` card) to create a sense of three-dimensional movement.

## 2. Colors & Surface Philosophy
The palette is a high-octane blend of the host nations' identities, filtered through a premium dark-mode lens.

*   **Primary (`#95aaff`):** Our "Electric Blue." Use this for key actions and interactive highlights.
*   **Secondary (`#ff706f`):** Our "Pulse Red." This provides the high-energy contrast essential for sports.
*   **Tertiary (`#dbffe9`):** Our "Pitch Green." Use this sparingly for success states or specialized "On-Pitch" data.

### The "No-Line" Rule
Standard 1px borders are strictly prohibited for sectioning. They clutter the UI and feel "cheap." Instead, define boundaries through **Background Color Shifts**.
*   Place a `surface-container-low` section on top of a `background` (`#0e0e0e`) to create a natural, sophisticated edge.
*   Use `surface-container-highest` (`#262626`) only for the most critical interactive components (like active match cards).

### Surface Hierarchy & Nesting
Treat the UI as a series of physical layers. 
1.  **Base:** `background` (#0e0e0e)
2.  **Sectioning:** `surface-container-low` (#131313)
3.  **Content Modules:** `surface-container` (#1a1a1a)
4.  **Interactive Elements:** `surface-bright` (#2c2c2c)

### The "Glass & Gradient" Rule
To elevate the system, use **Glassmorphism** for floating overlays (like live score updates or navigation bars). Utilize semi-transparent versions of `surface` tokens with a `backdrop-filter: blur(20px)`. 
*   **Signature Texture:** Main CTAs should use a linear gradient from `primary` (#95aaff) to `primary-container` (#829bff) at a 135-degree angle to provide "soul" and depth.

## 3. Typography
Our typography is the "Voice of the Commentator"—bold, clear, and varied.

*   **Display & Headlines (Lexend):** This is our "Athletic" face. Use `display-lg` (3.5rem) for hero scores and "Big Moment" headlines. The geometric nature of Lexend provides a modern, stadium-signage feel.
*   **Titles & Body (Manrope):** Manrope offers high legibility for dense match data. Use `title-lg` for team names and `body-md` for match statistics.
*   **Data Labels (Space Grotesk):** For technical data (e.g., "POSSESSION %", "xG"), use `label-md`. The monospaced feel of Space Grotesk adds a layer of "Technical Precision" to the broadcast aesthetic.

## 4. Elevation & Depth
We eschew traditional "drop shadows" in favor of **Tonal Layering**.

*   **The Layering Principle:** Depth is achieved by "stacking." A `surface-container-lowest` card placed on a `surface-container-low` section creates a recessed, "inground" look. Conversely, a `surface-bright` element on a `surface-container` creates a natural lift.
*   **Ambient Shadows:** For floating elements (Modals, Hovering Score Popups), use a shadow color tinted with `on-surface` at 5% opacity with a 40px blur. It should feel like a soft glow, not a hard shadow.
*   **The "Ghost Border" Fallback:** If containment is required for accessibility, use `outline-variant` at 15% opacity. Never use 100% opaque lines.

## 5. Components

### Buttons
*   **Primary:** Gradient (`primary` to `primary-container`), white text (`on-primary`), `md` (0.375rem) roundedness. 
*   **Secondary:** Ghost style. No background, `outline-variant` (20% opacity) border, `primary` text.
*   **Tertiary:** `surface-variant` background with `on-surface` text for low-priority actions.

### Cards (Match & Player)
*   **Strict Rule:** No dividers. Separate the "Home Team" and "Away Team" sections using a subtle shift from `surface-container` to `surface-container-high`.
*   **Layout:** Use asymmetrical padding (e.g., more padding on the left than the right) to create a sense of forward motion.

### Data Chips
*   Use `full` (9999px) roundedness. For "LIVE" indicators, use `secondary` (Pulse Red) with a subtle pulse animation. For "Full Time," use `surface-variant`.

### Inputs & Search
*   Use `surface-container-highest` as the base. On focus, do not use a heavy border; instead, transition the background to `surface-bright` and add a `primary` "glow" (ambient shadow).

### Specialized Sports Components
*   **The Score Streak:** A horizontal bar for match timelines. Use `primary` for team A and `secondary` for team B, set against a `surface-container-lowest` track.
*   **Momentum Tracker:** An area chart using a gradient of `primary` to transparent, showing which team is dominating the game.

## 6. Do's and Don'ts

### Do
*   **DO** use extreme scale. A `display-lg` score next to a `label-sm` timestamp creates high-end editorial drama.
*   **DO** lean into "Dark Mode" as the default. It mimics the premium feel of a night match under stadium lights.
*   **DO** use `secondary_dim` for "Alert" states to maintain the color story while ensuring readability.

### Don't
*   **DON'T** use 1px solid lines to separate content. Use the spacing scale and background shifts.
*   **DON'T** use standard "Web Blue." Stick strictly to our `primary` (#95aaff).
*   **DON'T** use `full` roundedness on anything except Chips and Avatars. We want the "Bold Athletic" feel of `md` (0.375rem) or `lg` (0.5rem) corners for structural elements.