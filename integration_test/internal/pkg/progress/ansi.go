package progress

import "fmt"

// esc generates multiple ANSI control characters
func esc(suffix ...string) (ansis string) {
	for _, s := range suffix {
		ansis += fmt.Sprintf("%c[%s", 033, s)
	}
	return
}

func clearLine() string {
	return esc("2K")
}

func cursorHorizontalAbsolute(n int) string {
	return esc(fmt.Sprintf("%dG", n))
}

func moveToLineHead() string {
	return cursorHorizontalAbsolute(1)
}

func cursorUp(n int) string {
	return esc(fmt.Sprintf("%dA", n))
}

// cursorUpHead up and move to head
func cursorUpHead(n int) string {
	return esc(fmt.Sprintf("%dF", n))
}

func saveCursor() string {
	return esc("s")
}

func loadCursor() string {
	return esc("u")
}
