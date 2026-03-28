package render

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

type barStyle struct{ filled, empty string }

var barStyles = map[string]barStyle{
	"blocks": {"█", "░"},
	"dots":   {"▮", "▯"},
	"ascii":  {"=", "-"},
	"thin":   {"━", "─"},
}

func ParseBarValue(s string) (float64, error) {
	if parts := strings.Split(s, "/"); len(parts) == 2 {
		num, err1 := strconv.ParseFloat(parts[0], 64)
		den, err2 := strconv.ParseFloat(parts[1], 64)
		if err1 != nil || err2 != nil || den == 0 {
			return 0, fmt.Errorf("invalid fraction %q", s)
		}
		return math.Round(num / den * 100), nil
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid value %q; expected number, decimal (0.0-1.0), or fraction (N/M)", s)
	}
	if v > 0 && v < 1 {
		return math.Round(v * 100), nil
	}
	if v == 1.0 {
		return 100, nil
	}
	return v, nil
}

func RenderBar(percent float64, width int, style string, showPercent bool) string {
	if percent < 0 { percent = 0 }
	if percent > 100 { percent = 100 }
	bs, ok := barStyles[style]
	if !ok { bs = barStyles["blocks"] }
	barWidth := width
	label := ""
	if showPercent {
		label = fmt.Sprintf(" %.0f%%", percent)
		barWidth = width - displayWidth(label)
	}
	if barWidth < 1 { barWidth = 1 }
	filled := int(math.Round(float64(barWidth) * percent / 100))
	empty := barWidth - filled
	return strings.Repeat(bs.filled, filled) + strings.Repeat(bs.empty, empty) + label
}
