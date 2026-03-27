package render

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

var sparkBars = []rune("▁▂▃▄▅▆▇█")

func ParseSparkValues(s string) ([]float64, error) {
	s = strings.TrimSpace(s)
	var parts []string
	if strings.Contains(s, ",") {
		parts = strings.Split(s, ",")
	} else {
		parts = strings.Fields(s)
	}
	values := make([]float64, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		v, err := strconv.ParseFloat(p, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid value %q: %w", p, err)
		}
		values = append(values, v)
	}
	if len(values) == 0 {
		return nil, fmt.Errorf("no values provided")
	}
	return values, nil
}

func RenderSpark(values []float64) string {
	if len(values) == 0 {
		return ""
	}
	min, max := values[0], values[0]
	for _, v := range values {
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}
	var sb strings.Builder
	for _, v := range values {
		idx := 0
		if max > min {
			idx = int(math.Round((v - min) / (max - min) * float64(len(sparkBars)-1)))
		}
		if idx < 0 {
			idx = 0
		}
		if idx >= len(sparkBars) {
			idx = len(sparkBars) - 1
		}
		sb.WriteRune(sparkBars[idx])
	}
	return sb.String()
}
