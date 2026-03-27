package main

import (
	"fmt"
	"os"
	"time"

	"golang.org/x/term"
)

func main() {
	w, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil || w <= 0 {
		w = 80
	}

	fmt.Fprintf(os.Stderr, "Terminal width: %d\n\n", w)

	// TEST 1: Can we overwrite a single line?
	fmt.Println("=== TEST 1: Single line overwrite ===")
	time.Sleep(300 * time.Millisecond)
	fmt.Fprint(os.Stdout, "Frame A")
	time.Sleep(300 * time.Millisecond)
	fmt.Fprint(os.Stdout, "\r\033[2KFrame B")
	time.Sleep(300 * time.Millisecond)
	fmt.Fprint(os.Stdout, "\r\033[2KFrame C\n")
	fmt.Println("(should see only 'Frame C' above)")
	time.Sleep(500 * time.Millisecond)

	// TEST 2: Can we overwrite 2 lines with cursor-up?
	fmt.Println("\n=== TEST 2: Two-line overwrite with cursor-up ===")
	time.Sleep(300 * time.Millisecond)
	fmt.Fprint(os.Stdout, "Line 1 - A\nLine 2 - A")
	time.Sleep(300 * time.Millisecond)
	// Move up 1 to Line 1, clear both, rewrite
	fmt.Fprint(os.Stdout, "\033[1A\r\033[2KLine 1 - B\n\033[2KLine 2 - B")
	time.Sleep(300 * time.Millisecond)
	fmt.Fprint(os.Stdout, "\033[1A\r\033[2KLine 1 - C\n\033[2KLine 2 - C\n")
	fmt.Println("(should see only 'Line 1 - C' and 'Line 2 - C' above)")
	time.Sleep(500 * time.Millisecond)

	// TEST 3: Same thing but with \033[0G instead of \r
	fmt.Println("\n=== TEST 3: Two-line overwrite with \\033[0G ===")
	time.Sleep(300 * time.Millisecond)
	fmt.Fprint(os.Stdout, "Line 1 - A\nLine 2 - A")
	time.Sleep(300 * time.Millisecond)
	fmt.Fprint(os.Stdout, "\033[1A\033[0G\033[2KLine 1 - B\n\033[0G\033[2KLine 2 - B")
	time.Sleep(300 * time.Millisecond)
	fmt.Fprint(os.Stdout, "\033[1A\033[0G\033[2KLine 1 - C\n\033[0G\033[2KLine 2 - C\n")
	fmt.Println("(should see only 'Line 1 - C' and 'Line 2 - C' above)")
	time.Sleep(500 * time.Millisecond)

	// TEST 4: Same with \033[1A repeated (BuildKit style)
	fmt.Println("\n=== TEST 4: BuildKit style (repeated \\033[1A) ===")
	time.Sleep(300 * time.Millisecond)
	fmt.Fprint(os.Stdout, "Line 1 - A\nLine 2 - A\nLine 3 - A")
	time.Sleep(300 * time.Millisecond)
	// Up 2 using repeated \033[1A
	fmt.Fprint(os.Stdout, "\033[1A\033[1A\033[0GLine 1 - B                    \n\033[0GLine 2 - B                    \n\033[0GLine 3 - B                    ")
	time.Sleep(300 * time.Millisecond)
	fmt.Fprint(os.Stdout, "\033[1A\033[1A\033[0GLine 1 - C                    \n\033[0GLine 2 - C                    \n\033[0GLine 3 - C                    \n")
	fmt.Println("(should see only C versions above)")
	time.Sleep(500 * time.Millisecond)

	// TEST 5: What our actual renderer does - space padding to full width
	fmt.Println("\n=== TEST 5: Space-padded to width (what gloss does) ===")
	time.Sleep(300 * time.Millisecond)
	pad := func(s string) string {
		p := w - 1 - len(s) // width-1 padding
		if p < 0 {
			p = 0
		}
		result := s
		for i := 0; i < p; i++ {
			result += " "
		}
		return result
	}
	fmt.Fprintf(os.Stdout, "\033[0G%s\n\033[0G%s\n", pad("--- gloss ---"), pad("  Frame A"))
	time.Sleep(300 * time.Millisecond)
	fmt.Fprintf(os.Stdout, "\033[1A\033[1A\033[0G%s\n\033[0G%s\n", pad("--- gloss ---"), pad("  Frame B"))
	time.Sleep(300 * time.Millisecond)
	fmt.Fprintf(os.Stdout, "\033[1A\033[1A\033[0G%s\n\033[0G%s\n", pad("--- gloss ---"), pad("  Frame C"))
	time.Sleep(300 * time.Millisecond)
	fmt.Fprintf(os.Stdout, "\033[1A\033[1A\033[0G%s\n\033[0G%s\n", pad("--- gloss ---"), pad("  Frame D"))
	time.Sleep(300 * time.Millisecond)
	fmt.Fprintf(os.Stdout, "\033[1A\033[1A\033[0G%s\n\033[0G%s\n", pad("--- gloss ---"), pad("  Frame E"))
	fmt.Println("(should see only ONE '--- gloss ---' and 'Frame E' above)")
	time.Sleep(500 * time.Millisecond)

	// TEST 6: What about with cursor hiding?
	fmt.Println("\n=== TEST 6: With hide/show cursor ===")
	time.Sleep(300 * time.Millisecond)
	fmt.Fprintf(os.Stdout, "\033[?25l\033[0G%s\n\033[0G%s\n\033[?25h", pad("--- gloss ---"), pad("  Frame A"))
	time.Sleep(300 * time.Millisecond)
	fmt.Fprintf(os.Stdout, "\033[?25l\033[1A\033[1A\033[0G%s\n\033[0G%s\n\033[?25h", pad("--- gloss ---"), pad("  Frame B"))
	time.Sleep(300 * time.Millisecond)
	fmt.Fprintf(os.Stdout, "\033[?25l\033[1A\033[1A\033[0G%s\n\033[0G%s\n\033[?25h", pad("--- gloss ---"), pad("  Frame C"))
	time.Sleep(300 * time.Millisecond)
	fmt.Fprintf(os.Stdout, "\033[?25l\033[1A\033[1A\033[0G%s\n\033[0G%s\n\033[?25h", pad("--- gloss ---"), pad("  Frame D"))
	fmt.Println("(should see only ONE '--- gloss ---' and 'Frame D' above)")

	fmt.Println("\n=== ALL TESTS DONE ===")
	fmt.Println("Which tests worked? Which showed duplicated lines?")
	fmt.Println("Please paste the output back.")
}
