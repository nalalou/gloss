package protocol

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/nalalou/gloss/internal/render"
	"golang.org/x/term"
)

func RenderLine(line string, width int, noColor bool) string {
	if !strings.HasPrefix(line, "::") {
		return line
	}
	rest := strings.TrimPrefix(line, "::")
	if rest == "" {
		return line
	}
	parts := strings.SplitN(rest, " ", 2)
	directive := strings.ToLower(parts[0])
	args := ""
	if len(parts) > 1 {
		args = parts[1]
	}

	switch directive {
	case "ok":
		return fmtBadge(args, "success", noColor)
	case "err", "error":
		return fmtBadge(args, "error", noColor)
	case "warn", "warning":
		return fmtBadge(args, "warning", noColor)
	case "info":
		return fmtBadge(args, "info", noColor)
	case "bar":
		return fmtBar(args, width)
	case "divider", "div":
		return render.RenderDivider(args, width, "light")
	case "callout":
		return fmtCallout(args)
	case "badge":
		return fmtExplicitBadge(args, noColor)
	case "list":
		return fmtList(args)
	case "table":
		return fmtTable(args)
	case "kv":
		return fmtKV(args)
	case "spark":
		return fmtSpark(args)
	case "color":
		return fmtColor(args, noColor)
	case "rain":
		return ""
	default:
		return line
	}
}

func fmtBadge(text string, badgeType string, noColor bool) string {
	result := render.RenderBadge(text, badgeType)
	if noColor {
		return result
	}
	color := render.BadgeDefaultColor(badgeType)
	if color != "" {
		return render.RenderStyled(result, color, false, false)
	}
	return result
}

func fmtBar(args string, width int) string {
	parts := strings.SplitN(strings.TrimSpace(args), " ", 2)
	if len(parts) == 0 || parts[0] == "" {
		return args
	}
	cleaned := strings.TrimSuffix(parts[0], "%")
	val, err := strconv.ParseFloat(cleaned, 64)
	if err != nil {
		return args
	}
	if val > 0 && val < 1 {
		val = val * 100
	} else if val == 1.0 && strings.Contains(parts[0], ".") {
		val = 100
	}
	label := ""
	if len(parts) > 1 {
		label = parts[1]
	}
	if width == 0 {
		width = termWidth()
	}
	result := render.RenderBar(val, width, "blocks", true)
	if label != "" {
		result = label + " " + result
	}
	return result
}

func fmtCallout(args string) string {
	parts := strings.SplitN(strings.TrimSpace(args), " ", 2)
	calloutType := "info"
	text := args
	if len(parts) >= 2 {
		switch parts[0] {
		case "success", "error", "warning", "info":
			calloutType = parts[0]
			text = parts[1]
		}
	}
	return render.RenderCallout(text, calloutType)
}

func fmtExplicitBadge(args string, noColor bool) string {
	parts := strings.SplitN(strings.TrimSpace(args), " ", 2)
	if len(parts) < 2 {
		return args
	}
	switch parts[0] {
	case "success", "error", "warning", "info", "plain":
		return fmtBadge(parts[1], parts[0], noColor)
	default:
		return fmtBadge(args, "plain", noColor)
	}
}

func fmtList(args string) string {
	items := strings.Split(args, "|")
	var parsed []string
	for _, item := range items {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		parts := strings.SplitN(item, " ", 2)
		if len(parts) == 2 {
			switch parts[0] {
			case "done", "fail", "pending":
				parsed = append(parsed, parts[1]+":"+parts[0])
				continue
			}
		}
		parsed = append(parsed, item)
	}
	if len(parsed) == 0 {
		return args
	}
	return render.RenderList(parsed, "bullet", true)
}

func fmtTable(args string) string {
	items := splitPipe(args)
	if len(items) == 0 {
		return args
	}
	rows := render.ParseKVArgs(items)
	return render.RenderTable(rows, "rounded")
}

func fmtKV(args string) string {
	items := splitPipe(args)
	if len(items) == 0 {
		return args
	}
	pairs := render.ParseKVPairs(items)
	return render.RenderKV(pairs, ":")
}

func fmtSpark(args string) string {
	values, err := render.ParseSparkValues(strings.TrimSpace(args))
	if err != nil {
		return args
	}
	return render.RenderSpark(values)
}

func fmtColor(args string, noColor bool) string {
	parts := strings.SplitN(strings.TrimSpace(args), " ", 2)
	if len(parts) < 2 {
		return args
	}
	if noColor {
		return parts[1]
	}
	hex, ok := render.ResolveColor(parts[0])
	if !ok {
		return args
	}
	return render.RenderStyled(parts[1], hex, false, false)
}

func splitPipe(s string) []string {
	items := strings.Split(s, "|")
	var result []string
	for _, item := range items {
		item = strings.TrimSpace(item)
		if item != "" {
			result = append(result, item)
		}
	}
	return result
}

func termWidth() int {
	w, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil || w <= 0 {
		return 80
	}
	return w
}

// Needed so render package functions used here are importable
var _ = fmt.Sprintf
