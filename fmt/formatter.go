package fmt

import (
	"fmt"
	"io"

	"github.com/robertgzr/color"
	"github.com/robertgzr/porcelain"
)

type Attrs = []color.Attribute

func attrs(a ...color.Attribute) Attrs {
	return a
}

type Formatter struct {
	NoColor                     bool
	ModeZsh, ModeBash, ModeTmux bool

	fgColor                  Attrs
	branchColor, commitColor Attrs
	aheadColor, behindColor  Attrs
	unmergedColor            Attrs
	untrackedColor           Attrs
	unstagedColor            Attrs
	dirtyColor, cleanColor   Attrs

	branchGlyph    string
	aheadGlyph     string
	behindGlyph    string
	unstagedGlyph  string
	untrackedGlyph string
	unmergedGlyph  string
	dirtyGlyph     string
	cleanGlyph     string
}

func (f *Formatter) Format(w io.Writer, repo porcelain.RepoState) {
	color.NoColor = f.NoColor
	color.EscapeBashPrompt = f.ModeBash
	color.EscapeZshPrompt = f.ModeZsh
	color.TmuxMode = f.ModeTmux

	color.New(f.fgColor...).Fprint(w, f.branchGlyph)
	fmt.Fprint(w, " ")
	color.New(f.branchColor...).Fprint(w, repo.Branch)

	color.New(f.fgColor...).Fprint(w, "@")
	color.New(f.commitColor...).Fprint(w, repo.Commit)
	fmt.Fprint(w, " ")

	if repo.CommitsAhead > 0 {
		color.New(f.aheadColor...).Fprint(w, repo.CommitsAhead)
	}
	if repo.CommitsBehind > 0 {
		color.New(f.behindColor...).Fprint(w, repo.CommitsBehind)
	}

	fmt.Fprint(w, " ")

	if repo.Untracked {
		color.New(f.untrackedColor...).Fprint(w, f.untrackedGlyph)
	}
	if repo.Unmerged {
		color.New(f.unmergedColor...).Fprint(w, f.unmergedGlyph)
	}
	if repo.Unstaged {
		color.New(f.unstagedColor...).Fprint(w, f.unstagedGlyph)
	}

	fmt.Fprint(w, " ")

	if repo.Clean {
		color.New(f.cleanColor...).Fprint(w, f.cleanGlyph)
	} else {
		color.New(f.dirtyColor...).Fprint(w, f.dirtyGlyph)
	}
}

func New() *Formatter {
	return &Formatter{
		NoColor:  false,
		ModeBash: false,
		ModeZsh:  false,
		ModeTmux: false,

		fgColor:        fgColor,
		branchColor:    branchColor,
		commitColor:    commitColor,
		aheadColor:     aheadColor,
		behindColor:    behindColor,
		unmergedColor:  unmergedColor,
		untrackedColor: untrackedColor,
		unstagedColor:  unstagedColor,
		dirtyColor:     dirtyColor,
		cleanColor:     cleanColor,

		branchGlyph:    branchGlyph,
		aheadGlyph:     aheadGlyph,
		behindGlyph:    behindGlyph,
		unstagedGlyph:  unstagedGlyph,
		untrackedGlyph: untrackedGlyph,
		unmergedGlyph:  unmergedGlyph,
		dirtyGlyph:     dirtyGlyph,
		cleanGlyph:     cleanGlyph,
	}
}
