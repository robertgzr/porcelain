package fmt

import (
	"github.com/robertgzr/color"
)

// see https://godoc.org/github.com/robertgzr/color#Attribute
var (
	fgColor     = attrs(color.FgWhite)
	branchColor = attrs(color.FgBlue)
	commitColor = attrs(color.FgGreen)

	aheadColor  = attrs(color.Faint, color.BgYellow, color.FgBlack)
	behindColor = attrs(color.Faint, color.BgRed, color.FgWhite)

	unstagedColor  = attrs(color.FgBlue)
	untrackedColor = attrs(color.Faint)
	unmergedColor  = attrs(color.FgCyan)

	dirtyColor = attrs(color.FgRed)
	cleanColor = attrs(color.FgGreen)
)

const (
	branchGlyph string = ""

	aheadGlyph  string = "↑"
	behindGlyph string = "↓"

	unstagedGlyph  string = "Δ"
	untrackedGlyph string = "?"
	unmergedGlyph  string = "‼"

	dirtyGlyph string = "✘"
	cleanGlyph string = "✔"
)
