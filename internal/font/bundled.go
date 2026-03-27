package font

import (
	"bytes"
	"embed"
	"fmt"
	"os"
	"strings"
)

//go:embed data/*.flf
var fontFS embed.FS

var bundledNames = map[string]string{
	"block":  "data/block.flf",
	"slant":  "data/slant.flf",
	"shadow": "data/shadow.flf",
	"doom":   "data/doom.flf",
	"small":  "data/small.flf",
	"big":    "data/big.flf",
	"thin":   "data/thin.flf",
	"3d":     "data/3d.flf",
	"script": "data/script.flf",
	"lean":   "data/lean.flf",
	"calvin": "data/calvin.flf",
	"banner": "data/banner.flf",
}

// BundledFontNames returns all available bundled font names.
func BundledFontNames() []string {
	return []string{"block", "slant", "shadow", "doom", "small", "big", "thin", "3d", "script", "lean", "calvin", "banner"}
}

// Load returns a Font by name (bundled) or by file path.
func Load(name string) (Font, error) {
	if strings.HasSuffix(name, ".flf") || strings.Contains(name, "/") || strings.Contains(name, "\\") {
		f, err := os.Open(name)
		if err != nil {
			return nil, fmt.Errorf("open font file %q: %w", name, err)
		}
		defer f.Close()
		return ParseFLF(f)
	}

	embeddedPath, ok := bundledNames[name]
	if !ok {
		return nil, fmt.Errorf("unknown font %q; run 'gloss fonts' to list available fonts", name)
	}
	data, err := fontFS.ReadFile(embeddedPath)
	if err != nil {
		return nil, fmt.Errorf("read bundled font %q: %w", name, err)
	}
	return ParseFLF(bytes.NewReader(data))
}
