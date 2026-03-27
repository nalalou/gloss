package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const glossTomlTemplate = `# gloss.toml — theme configuration for gloss
# All keys are optional. Flags override these values.
# Run 'gloss fonts' to preview available fonts.

# font = "block"
# gradient = ["#FF6B9D", "#6B9DFF"]
# shadow = false
# align = "left"
# border = "rounded"
# animate = false
`

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Scaffold a gloss.toml in the current directory",
	RunE:  runInit,
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func runInit(cmd *cobra.Command, args []string) error {
	path := "gloss.toml"

	if _, err := os.Stat(path); err == nil {
		return fmt.Errorf("gloss.toml already exists in current directory")
	}

	if err := os.WriteFile(path, []byte(glossTomlTemplate), 0644); err != nil {
		return fmt.Errorf("write gloss.toml: %w", err)
	}

	fmt.Printf("Created gloss.toml\nEdit it to set your theme, then run: gloss \"Your Text\"\n")
	return nil
}
