package protocol

import (
	"strings"

	"github.com/nalalou/gloss/internal/render"
)

// GitHub Actions workflow commands: "::name::message" or
// "::name key=value,key=value::message".
var githubCommands = map[string]bool{
	"error":    true,
	"warning":  true,
	"notice":   true,
	"debug":    true,
	"group":    true,
	"endgroup": true,
	"add-mask": true,
}

// ParseGitHub recognizes a GitHub Actions workflow command.
// ok is false for anything else, including gloss-style lines like "::warn Slow::query".
func ParseGitHub(line string) (name string, params map[string]string, message string, ok bool) {
	if !strings.HasPrefix(line, "::") {
		return "", nil, "", false
	}
	rest := line[2:]
	end := strings.Index(rest, "::")
	if end < 0 {
		return "", nil, "", false
	}
	head := rest[:end]
	message = rest[end+2:]

	name, paramStr, _ := strings.Cut(head, " ")
	name = strings.ToLower(name)
	if !githubCommands[name] {
		return "", nil, "", false
	}

	params = map[string]string{}
	if paramStr != "" {
		for _, kv := range strings.Split(paramStr, ",") {
			k, v, found := strings.Cut(strings.TrimSpace(kv), "=")
			if !found || k == "" || strings.Contains(k, " ") {
				// Not key=value pairs, so this is a gloss-style line, e.g. "::warning Slow query::x".
				return "", nil, "", false
			}
			params[k] = unescapeGitHub(v)
		}
	}
	return name, params, unescapeGitHub(message), true
}

// unescapeGitHub reverses the percent-encoding GitHub uses for data and properties.
func unescapeGitHub(s string) string {
	r := strings.NewReplacer("%0D", "\r", "%0A", "\n", "%3A", ":", "%2C", ",", "%25", "%")
	return r.Replace(s)
}

func renderGitHub(name string, params map[string]string, message string, width int, noColor bool) string {
	switch name {
	case "error", "warning", "notice":
		text := message
		if title := params["title"]; title != "" {
			text = title + ": " + text
		}
		if loc := githubLocation(params); loc != "" {
			text = loc + " " + text
		}
		badgeType := map[string]string{"error": "error", "warning": "warning", "notice": "info"}[name]
		return fmtBadge(text, badgeType, noColor)
	case "group":
		return render.RenderDivider(message, width, "light")
	default: // endgroup, debug, add-mask: GitHub shows nothing for these
		return ""
	}
}

func githubLocation(params map[string]string) string {
	file := params["file"]
	if file == "" {
		return ""
	}
	if line := params["line"]; line != "" {
		file += ":" + line
		if col := params["col"]; col != "" {
			file += ":" + col
		}
	}
	return file
}
