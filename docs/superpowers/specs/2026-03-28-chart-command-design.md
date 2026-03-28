# `gloss chart` command

## Summary

A new subcommand that renders tall vertical bar charts in the terminal with gradient coloring and optional animation. Accepts comma-separated numeric values via positional args or stdin.

## CLI Interface

```bash
gloss chart 220,195,180,140,95,72,45
echo "220,195,180,140,95" | gloss chart
gloss chart 220,195,180 --height=8 --label="p99 (ms)" --gradient=fire
```

## Flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--height` | int | 6 | Number of rows tall |
| `--label` | string | (none) | Dim label printed above the chart |
| `--gradient` | string | (inherited) | PersistentFlag from root |
| `--color` | string | (inherited) | PersistentFlag from root |
| `--no-color` | bool | (inherited) | PersistentFlag from root |
| `--width` | int | (inherited) | PersistentFlag from root |

## Rendering Rules

- Each bar is 2 chars wide (`██`) with 1 char gap between bars.
- Half-block `▄` used for smooth tops when a bar height falls between integer rows.
- Gradient is interpolated across columns left-to-right.
- 2-space left indent for visual breathing room.
- Values are scaled relative to the maximum value in the dataset (max value = full height).

## Animation

- **TTY detected:** Columns appear one at a time. Base delay 120ms, +10ms per column for slight acceleration feel. Cursor hidden during animation, restored on exit/interrupt.
- **Non-TTY (piped):** Full chart printed at once, no animation, no cursor manipulation.

## Architecture

### `internal/render/chart.go`

Pure rendering function, no animation, no color, no terminal awareness.

```go
func RenderChart(values []float64, height int) []string
```

- Takes values and height, returns slice of strings (one per row, top to bottom).
- Each bar is 2 chars wide, 1 char gap. Uses `█` for full cells, `▄` for half cells, space for empty.
- 2-space left indent on each line.
- Does not handle color or gradient — that's the cmd layer's job via `ColorizeLines`.

### `internal/render/chart_test.go`

Test cases:
- Basic chart with varied values
- Single value — one bar at full height
- All same values — all bars at full height
- All zeros — flat line of `▄` at bottom row
- Two values — minimum viable chart
- Negative values — clamped to 0

### `cmd/chart.go`

Cobra command. Responsibilities:
- Parse comma-separated values from positional args or stdin (using `readStdinText` helper from `cmd/helpers.go`)
- Call `render.RenderChart` for the lines
- Apply `render.ColorizeLines` with gradient/color flags
- If TTY: animate column-by-column with cursor movement. Reuse `animateTallBars` logic from demo.go (move shared animation code here, have demo.go call into it or duplicate the small animation loop).
- If non-TTY: print all lines at once
- Print dim label above chart if `--label` is set
- Error on empty input: `"no data provided; see 'gloss chart --help'"`
- Skip non-numeric values with warning to stderr

### `cmd/demo.go` update

Replace the inline `animateTallBars` call with the new `chart` rendering. The demo can either:
- Call `render.RenderChart` directly and animate inline (keeping demo self-contained)
- Or import the animation helper if extracted to a shared location

Recommendation: keep the animation helpers in demo.go since they're demo-specific (spinner, typewriter, etc). The chart command gets its own animation loop in cmd/chart.go. Some duplication is fine — the animation loops are ~15 lines each and serve different contexts.

## Edge Cases

| Input | Behavior |
|-------|----------|
| No input | Error: "no data provided; see 'gloss chart --help'" |
| Single value | One bar at full height |
| All zeros | Flat line of `▄` at bottom |
| All same values | All bars at full height |
| Negative values | Clamped to 0 |
| Non-numeric values | Skipped, warning to stderr |
| Very many values (>50) | Render all, may exceed terminal width — no wrapping |

## Non-goals

- Axis labels or tick marks
- Horizontal bar charts
- Line charts or other chart types (future work)
- Interactive features
