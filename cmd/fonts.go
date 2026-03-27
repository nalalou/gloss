package cmd

import (
	"fmt"
	"strings"

	"github.com/nalalou/gloss/internal/font"
	"github.com/nalalou/gloss/internal/render"
	"github.com/nalalou/gloss/internal/theme"
	"github.com/spf13/cobra"
)

var fontsCmd = &cobra.Command{
	Use:   "fonts",
	Short: "Preview all bundled fonts",
	Long:  "Renders the word 'Gloss' in each bundled font so you can choose your favorite.",
	RunE:  runFonts,
}

func init() {
	rootCmd.AddCommand(fontsCmd)
}

func runFonts(cmd *cobra.Command, args []string) error {
	names := font.BundledFontNames()
	previewText := "Gloss"

	for _, name := range names {
		f, err := font.Load(name)
		if err != nil {
			fmt.Printf("  [%s: failed to load — %v]\n\n", name, err)
			continue
		}

		header := strings.Repeat("─", 40)
		fmt.Println(header)
		fmt.Printf("  font: %s\n", name)
		fmt.Println(header)

		opts := theme.Defaults()
		opts.Gradient = []string{"#FF6B9D", "#6B9DFF"}
		opts.Align = "left"

		output := render.Render(previewText, f, opts)
		fmt.Println(output)
		fmt.Println()
	}
	return nil
}
