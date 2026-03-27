package watch

import (
	"strings"

	"github.com/nalalou/gloss/internal/protocol"
	"github.com/nalalou/gloss/internal/render"
)

type Element struct {
	ID        string
	Directive string
	Args      string
	State     string
	Rendered  string
}

type Panel struct {
	order    []string
	elements map[string]*Element
	width    int
}

func NewPanel(width int) *Panel {
	return &Panel{
		elements: make(map[string]*Element),
		width:    width,
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
	if directive == "status" || directive == "spin" {
		parts := strings.SplitN(args, " ", 2)
		if len(parts) >= 1 {
			elem.State = parts[0]
		}
	}
	elem.Rendered = p.renderElement(elem, noColor)
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
	lines = append(lines, render.RenderDivider("gloss", p.width, "light"))
	for _, id := range p.order {
		if elem, ok := p.elements[id]; ok {
			for _, line := range strings.Split(elem.Rendered, "\n") {
				lines = append(lines, "  "+line)
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
			parts := strings.SplitN(elem.Args, " ", 2)
			text := ""
			if len(parts) >= 2 {
				text = parts[1]
			}
			elem.Rendered = icon + " " + text
		}
	}
}

func (p *Panel) SetWidth(width int) { p.width = width }

func (p *Panel) renderElement(elem *Element, noColor bool) string {
	line := "::" + elem.Directive + " " + elem.Args
	return protocol.RenderLine(line, p.width-4, noColor)
}
