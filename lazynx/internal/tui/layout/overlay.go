package layout

import (
	"bytes"
	"strings"

	"github.com/mattn/go-runewidth"
	"github.com/muesli/reflow/ansi"
	"github.com/muesli/reflow/truncate"
)

// WhitespaceOption configures whitespace rendering (for future extensibility)
type WhitespaceOption func(*whitespace)

type whitespace struct {
	chars string
}

func (w *whitespace) render(width int) string {
	if width <= 0 {
		return ""
	}
	if w.chars == "" {
		return strings.Repeat(" ", width)
	}
	return strings.Repeat(w.chars, width)
}

func clamp(v, low, high int) int {
	if v < low {
		return low
	}
	if v > high {
		return high
	}
	return v
}

func getLines(s string) ([]string, int) {
	lines := strings.Split(s, "\n")
	width := 0
	for _, line := range lines {
		lineWidth := ansi.PrintableRuneWidth(line)
		if lineWidth > width {
			width = lineWidth
		}
	}
	return lines, width
}

// PlaceOverlay places fg on top of bg.
// took from https://github.com/charmbracelet/lipgloss/pull/102/files
func PlaceOverlay(x, y int, fg, bg string, opts ...WhitespaceOption) string {
	fgLines, fgWidth := getLines(fg)
	bgLines, bgWidth := getLines(bg)
	bgHeight := len(bgLines)
	fgHeight := len(fgLines)

	if fgWidth >= bgWidth && fgHeight >= bgHeight {
		// If foreground is larger than background, just return foreground
		return fg
	}

	// Clamp position to ensure overlay fits within background
	x = clamp(x, 0, bgWidth-fgWidth)
	y = clamp(y, 0, bgHeight-fgHeight)

	ws := &whitespace{}
	for _, opt := range opts {
		opt(ws)
	}

	var b strings.Builder
	for i, bgLine := range bgLines {
		if i > 0 {
			b.WriteByte('\n')
		}
		if i < y || i >= y+fgHeight {
			b.WriteString(bgLine)
			continue
		}

		pos := 0
		if x > 0 {
			left := truncate.String(bgLine, uint(x))
			pos = ansi.PrintableRuneWidth(left)
			b.WriteString(left)
			if pos < x {
				b.WriteString(ws.render(x - pos))
				pos = x
			}
		}

		fgLine := fgLines[i-y]
		b.WriteString(fgLine)
		pos += ansi.PrintableRuneWidth(fgLine)

		right := cutLeft(bgLine, pos)
		bgWidth := ansi.PrintableRuneWidth(bgLine)
		rightWidth := ansi.PrintableRuneWidth(right)
		if rightWidth <= bgWidth-pos {
			b.WriteString(ws.render(bgWidth - rightWidth - pos))
		}

		b.WriteString(right)
	}

	return b.String()
}

// cutLeft cuts printable characters from the left.
// This function is heavily based on muesli's ansi and truncate packages.
// took from https://github.com/charmbracelet/lipgloss/pull/102/files
func cutLeft(s string, cutWidth int) string {
	var (
		pos    int
		isAnsi bool
		ab     bytes.Buffer
		b      bytes.Buffer
	)
	for _, c := range s {
		var w int
		if c == ansi.Marker || isAnsi {
			isAnsi = true
			ab.WriteRune(c)
			if ansi.IsTerminator(c) {
				isAnsi = false
				if bytes.HasSuffix(ab.Bytes(), []byte("[0m")) {
					ab.Reset()
				}
			}
		} else {
			w = runewidth.RuneWidth(c)
		}

		if pos >= cutWidth {
			if b.Len() == 0 {
				if ab.Len() > 0 {
					b.Write(ab.Bytes())
				}
				if pos-cutWidth > 1 {
					b.WriteByte(' ')
					continue
				}
			}
			b.WriteRune(c)
		}
		pos += w
	}
	return b.String()
}
