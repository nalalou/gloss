package cmd

import (
	"encoding/csv"
	"fmt"
	"strings"

	"github.com/nalalou/gloss/internal/env"
	"github.com/nalalou/gloss/internal/render"
	"github.com/spf13/cobra"
)

var (
	flagTableCSV   bool
	flagTableTSV   bool
	flagTableStyle string
)

var tableCmd = &cobra.Command{
	Use:   "table [key=value...]",
	Short: "Render a formatted table",
	Long: `Renders a bordered table. Pass key=value arguments for a two-column table,
or pipe CSV/TSV data with --csv or --tsv for arbitrary columns. The first
row of CSV/TSV input is treated as the header.`,
	Example: `  gloss table "Name=Alice" "Role=Engineer"
  echo "Name,Role\nAlice,Engineer" | gloss table --csv`,
	Args: cobra.ArbitraryArgs,
	RunE: runTable,
}

func init() {
	tableCmd.Flags().BoolVar(&flagTableCSV, "csv", false, "read CSV from stdin")
	tableCmd.Flags().BoolVar(&flagTableTSV, "tsv", false, "read TSV from stdin")
	tableCmd.Flags().StringVar(&flagTableStyle, "style", "rounded", "border style: none, single, double, rounded, thick")
	rootCmd.AddCommand(tableCmd)
}

func runTable(cmd *cobra.Command, args []string) error {
	var rows [][]string
	if flagTableCSV || flagTableTSV {
		raw, err := readStdinText(int64(maxInputSize))
		if err != nil {
			return err
		}
		r := csv.NewReader(strings.NewReader(raw))
		if flagTableTSV {
			r.Comma = '\t'
		}
		rows, err = r.ReadAll()
		if err != nil {
			return fmt.Errorf("parse CSV: %w", err)
		}
	} else if len(args) > 0 {
		rows = render.ParseKVArgs(args)
	} else {
		return fmt.Errorf("no data provided; see 'gloss table --help'")
	}
	if len(rows) == 0 {
		return fmt.Errorf("no data provided; see 'gloss table --help'")
	}

	result := render.RenderTable(rows, flagTableStyle)
	envInfo := env.Detect()
	noColor := envInfo.NoColor || flagNoColor
	gradient := resolveGradientFlag()
	lines := strings.Split(result, "\n")
	lines = render.ColorizeLines(lines, gradient, flagColor, noColor)
	fmt.Println(strings.Join(lines, "\n"))
	return nil
}
