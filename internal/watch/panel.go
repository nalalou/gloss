package watch

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/x/ansi"
	"github.com/nalalou/gloss/internal/protocol"
	"github.com/nalalou/gloss/internal/render"
)

type Element struct {
	ID        string
	Directive string
	Args      string
	State     string
	Rendered  string
	Started   time.Time     // when the element entered the running state
	Elapsed   time.Duration // how long it ran, once finished
}

type Panel struct {
	order    []string
	elements map[string]*Element
	width    int
	now      func() time.Time
}

func NewPanel(width int) *Panel {
	return &Panel{
		elements: make(map[string]*Element),
		width:    width,
		now:      time.Now,
	}
}

func (p *Panel) Set(id, directive, args string, noColor bool) {
	elem, exists := p.elements[id]
	if !exists {
		elem = &Element{ID: id}
		p.elements[id] = elem
		p.order = append(p.order, id)
	}
	elem.Directive = directive
	elem.Args = args
	prevState := elem.State
	elem.State = ""
	if directive == "status" {
		elem.State, _, _ = strings.Cut(args, " ")
	} else if directive == "spin" {
		elem.State = "running"
	}
	switch {
	case elem.State == "running" && prevState != "running":
		elem.Started = p.now()
		elem.Elapsed = 0
	case elem.State != "running" && prevState == "running":
		elem.Elapsed = p.now().Sub(elem.Started)
	}
	elem.Rendered = p.renderElement(elem, noColor)
	if elem.Elapsed > 0 {
		elem.Rendered += timing(elem.Elapsed, noColor)
	}
}

// timing renders a step's duration, or nothing if it was too quick to matter.
func timing(d time.Duration, noColor bool) string {
	if d < time.Second {
		return ""
	}
	text := " " + FormatDuration(d)
	if noColor {
		return text
	}
	return render.RenderStyled(text, "#888888", false, true)
}

// FormatDuration renders d as "4.2s", "38s", or "2m05s".
func FormatDuration(d time.Duration) string {
	switch {
	case d < 10*time.Second:
		return fmt.Sprintf("%.1fs", d.Seconds())
	case d < time.Minute:
		return fmt.Sprintf("%ds", int(d.Seconds()))
	default:
		return fmt.Sprintf("%dm%02ds", int(d.Minutes()), int(d.Seconds())%60)
	}
}

func (p *Panel) Remove(id string) {
	delete(p.elements, id)
	for i, oid := range p.order {
		if oid == id {
			p.order = append(p.order[:i], p.order[i+1:]...)
			break
		}
	}
}

func (p *Panel) Len() int { return len(p.order) }

func (p *Panel) HasRunning() bool {
	for _, id := range p.order {
		if elem, ok := p.elements[id]; ok && elem.State == "running" {
			return true
		}
	}
	return false
}

func (p *Panel) Height() int {
	if len(p.order) == 0 {
		return 0
	}
	h := 1 // divider
	for _, id := range p.order {
		if elem, ok := p.elements[id]; ok {
			h += strings.Count(elem.Rendered, "\n") + 1
		}
	}
	return h
}

func (p *Panel) RenderLines() []string {
	if len(p.order) == 0 {
		return nil
	}
	var lines []string
	lines = append(lines, render.RenderDivider("gloss", p.width-1, "light"))
	for _, id := range p.order {
		if elem, ok := p.elements[id]; ok {
			for _, line := range strings.Split(elem.Rendered, "\n") {
				// A line that wraps would throw off the cursor math that redraws the panel.
				lines = append(lines, ansi.Truncate("  "+line, p.width-1, "…"))
			}
		}
	}
	return lines
}

func (p *Panel) UpdateSpinnerFrame(frame int, noColor bool) {
	frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	for _, id := range p.order {
		elem := p.elements[id]
		if elem.State == "running" {
			icon := frames[frame%len(frames)]
			text := elem.Args
			if elem.Directive == "status" {
				_, text, _ = strings.Cut(elem.Args, " ")
			}
			elem.Rendered = icon + " " + text + timing(p.now().Sub(elem.Started), noColor)
		}
	}
}

func (p *Panel) SetWidth(width int) { p.width = width }

func (p *Panel) renderElement(elem *Element, noColor bool) string {
	line := "::" + elem.Directive + " " + elem.Args
	return protocol.RenderLine(line, p.width-4, noColor)
}
