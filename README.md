```
             __
   ____ _   / /  ____    _____   _____
  / __ `/  / /  / __ \  / ___/  / ___/
 / /_/ /  / /  / /_/ / (__  )  (__  )
 \__, /  /_/   \____/ /____/  /____/
/____/
```

**GitHub Actions** **`::`** **annotations, but for your terminal.**

Your scripts emit `::` directives. Gloss renders them as progress bars, badges, tables, sparklines, and live-updating panels in one binary you pipe through. Zero dependencies in your code, just `echo`.

```bash
brew install gloss  # coming soon
go install github.com/nalalou/gloss@latest
```

***

## The Idea

AI agents and build scripts produce walls of text. Gloss adds structure by piping through a formatter.

```bash
# Your agent/script just prints lines:
echo "::status id=build running Building..."
echo "compiling main.go..."
echo "::status id=build done Build complete"
echo "::bar id=prog 100 Progress"
echo "::ok All green"

# Pipe through gloss:
./my-script.sh | gloss watch
```

The `::` lines become live-updating widgets. Everything else scrolls normally. Works from Bash, Python, Go, Rust, TypeScript — anything that can `print`.

***

## What It Looks Like

```
compiling main.go...
compiling utils.go...
PASS auth_test.go (12 tests)
PASS api_test.go (8 tests)
───────────────────────────── gloss ─────────────────────────────
  ✓ Build
  ✓ Test (26 passed)
  ⠹ Deploy
  ████████████████████████████░░░░░░░░░░░ 66%
```

The bottom panel stays in place and updates live. The top scrolls.

***

## The Protocol

Lines starting with `::` are directives. Everything else passes through.

```
::ok Tests passing              ✓ Tests passing
::err Build failed              ✗ Build failed
::warn Slow query               ⚠ Slow query
::info Using cache              ℹ Using cache
::bar 75 Coverage               ████████████████████░░░░░░ 75%
::divider Results               ──────── Results ────────
::callout error Timeout         ╭ ✗ Error ──────╮
                                │ Timeout       │
                                ╰───────────────╯
```

Add `id=` to make it **live-updating** (tracked in the panel):

```
::status id=build running Building     ⠹ Building  (animates)
::status id=build done Build           ✓ Build     (replaces above)
::bar id=prog 75 Progress              updates in place
::kv id=meta Pods=3/3 | Region=us-1   updates in place
```

**No** **`id=`** → scrolls once. **With** **`id=`** → persists in panel, updates in place.

***

## Two Modes

### `gloss fmt` — stateless pipe formatter

Renders `::` directives as styled text. No cursor tricks. Pipe-safe, file-safe.

```bash
my-script.sh | gloss fmt
```

### `gloss watch` — live TUI panel

Splits terminal into scroll zone + persistent panel. `id=` directives update in place. Spinners animate. Falls back to `gloss fmt` when piped.

```bash
my-agent | gloss watch
```

***

## All the Primitives

Gloss also works as standalone subcommands for shell scripts:

```bash
# Badges
gloss badge "Tests passing" --type=success     # ✓ Tests passing
gloss badge "Build failed" --type=error        # ✗ Build failed

# Progress bars
gloss bar 75                                    # ████████████████░░░░ 75%
gloss bar 3/4 --style=dots                     # ▮▮▮▮▮▮▯▯▯▯ 75%

# Sparklines
gloss spark 1,4,7,3,9,2,5                     # ▁▃▆▂█▁▄

# Dividers
gloss divider "Section Title"                  # ──── Section Title ────

# Tables
gloss table "Name=Alice" "Role=Engineer"       # bordered table

# Key-value pairs
gloss kv "CPU=58%" "Mem=76%" "Disk=29%"        # aligned columns

# Lists with status
gloss list "Build:done" "Test:fail" --status   # ✓ Build  ✗ Test

# Callout boxes
gloss callout "Approval needed" --type=warning # boxed alert

# Inline color
echo "Status: $(gloss color OK --fg=green)"    # colored inline text

# Spinners
gloss spin "Installing..." -- npm install      # ⠹ Installing... → ✓

# ASCII art banners
gloss "Hello World" --font=doom --gradient=fire
```

***

## 12 Fonts

```bash
gloss fonts  # preview all
```

`block` · `slant` · `shadow` · `doom` · `small` · `big` · `thin` · `3d` · `script` · `lean` · `calvin` · `banner`

## 15 Gradients

`fire` · `ocean` · `mono` · `neon` · `aurora` · `sunset` · `synthwave` · `matrix` · `cyberpunk` · `pastel` · `lavender` · `ice` · `autumn` · `mint` · `rainbow`

Custom gradients with 2+ hex colors:

```bash
gloss "RGB" --gradient="#FF0000,#00FF00,#0000FF"
```

***

## For AI Agent Builders

Gloss was designed for agents. The `::` protocol is what GitHub Actions uses for `::warning`, `::group`, `::error` — but portable to any terminal.

An agent just prints lines. It doesn't import a library. It doesn't adopt a framework. It just `echo`s.

```python
# Python agent — no gloss dependency
print("::status id=research running Searching papers...")
print("Found 14 relevant papers")
print("::status id=research done Research complete")
print("::bar id=progress 50")
```

```bash
# Bash script — no gloss dependency
echo "::status id=deploy running Deploying..."
kubectl rollout status deployment/app
echo "::status id=deploy done Deploy complete"
```

```go
// Go tool — no gloss dependency
fmt.Println("::ok Build passed")
fmt.Println("::bar 100 Coverage")
```

The user runs: `my-agent | gloss watch`

***

## Config

```bash
gloss init  # creates gloss.toml
```

```toml
# gloss.toml
font = "block"
gradient = ["#FF6B9D", "#6B9DFF"]
border = "rounded"
# shadow = false
# align = "left"
# animate = false
```

Flags override config. `NO_COLOR` and `TERM=dumb` are respected. Color auto-disables when piped.

***

## Install

```bash
# Go
go install github.com/nalalou/gloss@latest

# From source
git clone https://github.com/nalalou/gloss.git
cd gloss && go build -o gloss .
```

Homebrew tap coming with first tagged release.

***

## The Taste Line

> Format the structure, not the content. Dividers mark sections. Badges mark outcomes. Callouts mark exceptions. Everything else is plain text.

If you're formatting more than 10% of your output, you've crossed the line.

***

## License

MIT

## Acknowledgments

Built with [Lip Gloss](https://github.com/charmbracelet/lipgloss) and inspired by the [Charmbracelet](https://charm.sh) ecosystem.
