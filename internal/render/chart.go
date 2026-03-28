package render

import "math"

func RenderChart(values []float64, height int) []string {
	if len(values) == 0 {
		return nil
	}

	clamped := make([]float64, len(values))
	for i, v := range values {
		if v < 0 {
			v = 0
		}
		clamped[i] = v
	}

	max := 0.0
	for _, v := range clamped {
		if v > max {
			max = v
		}
	}
	if max == 0 {
		max = 1
	}

	barWidth := 2
	gap := 1
	colWidth := barWidth + gap
	totalWidth := len(clamped)*colWidth - gap

	grid := make([][]rune, height)
	for r := range grid {
		grid[r] = make([]rune, totalWidth)
		for c := range grid[r] {
			grid[r][c] = ' '
		}
	}

	for col, v := range clamped {
		barH := v / max * float64(height)
		fullCells := int(barH)
		hasHalf := (barH - float64(fullCells)) >= 0.5

		if v == 0 && max > 0 {
			hasHalf = true
			fullCells = 0
		}

		x := col * colWidth
		for r := 0; r < height; r++ {
			rowFromBottom := height - 1 - r
			for b := 0; b < barWidth; b++ {
				if x+b < totalWidth {
					if rowFromBottom < fullCells {
						grid[r][x+b] = '█'
					} else if rowFromBottom == fullCells && hasHalf {
						grid[r][x+b] = '▄'
					}
				}
			}
		}
	}

	lines := make([]string, height)
	for r := 0; r < height; r++ {
		lines[r] = "  " + string(grid[r])
	}
	return lines
}

func ChartColumnCount(n int) int {
	if n == 0 {
		return 0
	}
	return n*3 - 1
}

func RenderChartPartial(values []float64, height int, numCols int) []string {
	if numCols <= 0 || numCols > len(values) {
		numCols = len(values)
	}
	max := 0.0
	for _, v := range values {
		if v > 0 && v > max {
			max = v
		}
	}
	if max == 0 {
		max = 1
	}

	subset := values[:numCols]
	barWidth := 2
	gap := 1
	colWidth := barWidth + gap
	fullWidth := len(values)*colWidth - gap

	grid := make([][]rune, height)
	for r := range grid {
		grid[r] = make([]rune, fullWidth)
		for c := range grid[r] {
			grid[r][c] = ' '
		}
	}

	for col, v := range subset {
		if v < 0 {
			v = 0
		}
		barH := v / max * float64(height)
		fullCells := int(barH)
		hasHalf := (barH - math.Floor(barH)) >= 0.5

		if v == 0 {
			hasHalf = true
			fullCells = 0
		}

		x := col * colWidth
		for r := 0; r < height; r++ {
			rowFromBottom := height - 1 - r
			for b := 0; b < barWidth; b++ {
				if x+b < fullWidth {
					if rowFromBottom < fullCells {
						grid[r][x+b] = '█'
					} else if rowFromBottom == fullCells && hasHalf {
						grid[r][x+b] = '▄'
					}
				}
			}
		}
	}

	lines := make([]string, height)
	for r := 0; r < height; r++ {
		lines[r] = "  " + string(grid[r])
	}
	return lines
}
