package porcelain

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/robertgzr/color"
)

type GitArea struct {
	modified int
	added    int
	deleted  int
	renamed  int
	copied   int
}

func (a *GitArea) HasChanged() bool {
	var changed bool
	if a.added != 0 {
		changed = true
	}
	if a.deleted != 0 {
		changed = true
	}
	if a.modified != 0 {
		changed = true
	}
	if a.copied != 0 {
		changed = true
	}
	if a.renamed != 0 {
		changed = true
	}
	return changed
}

type PorcInfo struct {
	workingDir string

	branch   string
	commit   string
	remote   string
	upstream string
	ahead    int
	behind   int

	untracked int
	unmerged  int

	Unstaged GitArea
	Staged   GitArea
}

func (pi *PorcInfo) HasUnmerged(cwd string) bool {
	if pi.unmerged > 0 {
		return true
	}
	gitDir, err := PathToGitDir(cwd)
	if err != nil {
		log.Printf("error calling PathToGitDir: %s", err)
		return false
	}
	// TODO figure out if output of MERGE_HEAD can be useful
	if _, err := ioutil.ReadFile(path.Join(gitDir, "MERGE_HEAD")); err != nil {
		if os.IsNotExist(err) {
			return false
		}
		log.Printf("error reading MERGE_HEAD: %s", err)
		return false
	} else {
		return true
	}
}
func (pi *PorcInfo) HasModified() bool {
	return pi.Unstaged.HasChanged()
}
func (pi *PorcInfo) IsDirty() bool {
	return pi.Staged.HasChanged()
}

func (pi *PorcInfo) Debug() string {
	return fmt.Sprintf("%#+v", pi)
}

// Fmt formats the output for the shell
// TODO should be configurable by the user
//
func (pi *PorcInfo) Fmt(cwd string, noColorFlag, bashFmtFlag, zshFmtFlag, tmuxFmtFlag bool) string {
	log.Printf("formatting output: %s", pi.Debug())

	var (
		branchGlyph   string = ""
		modifiedGlyph string = "Δ"
		// deletedGlyph   string = "＊"
		dirtyGlyph     string = "✘"
		cleanGlyph     string = "✔"
		untrackedGlyph string = "?"
		unmergedGlyph  string = "‼"
		aheadArrow     string = "↑"
		behindArrow    string = "↓"
	)

	if noColorFlag {
		color.NoColor = true
	} else {
		color.NoColor = false
		color.EscapeBashPrompt = bashFmtFlag
		color.EscapeZshPrompt = zshFmtFlag
		color.TmuxMode = tmuxFmtFlag
	}
	branchFmt := color.New(color.FgBlue).SprintFunc()
	commitFmt := color.New(color.FgGreen, color.Italic).SprintFunc()

	aheadFmt := color.New(color.Faint, color.BgYellow, color.FgBlack).SprintFunc()
	behindFmt := color.New(color.Faint, color.BgRed, color.FgWhite).SprintFunc()

	modifiedFmt := color.New(color.FgBlue).SprintFunc()
	// deletedFmt := color.New(color.FgYellow).SprintFunc()
	dirtyFmt := color.New(color.FgRed).SprintFunc()
	cleanFmt := color.New(color.FgGreen).SprintFunc()

	untrackedFmt := color.New(color.Faint).SprintFunc()
	unmergedFmt := color.New(color.FgCyan).SprintFunc()

	return fmt.Sprintf("%s %s@%s %s %s %s",
		branchGlyph,
		branchFmt(pi.branch),
		func() string {
			if pi.commit == "(initial)" {
				return commitFmt(pi.commit)
			}
			return commitFmt(pi.commit[:7])
		}(),
		func() string {
			var buf bytes.Buffer
			if pi.ahead > 0 {
				buf.WriteString(aheadFmt(" ", aheadArrow, pi.ahead, " "))
			}
			if pi.behind > 0 {
				buf.WriteString(behindFmt(" ", behindArrow, pi.behind, " "))
			}
			return buf.String()
		}(),
		func() string {
			var buf bytes.Buffer
			if pi.untracked > 0 {
				buf.WriteString(untrackedFmt(untrackedGlyph))
			} else {
				buf.WriteRune(' ')
			}
			if pi.HasUnmerged(cwd) {
				buf.WriteString(unmergedFmt(unmergedGlyph))
			} else {
				buf.WriteRune(' ')
			}
			if pi.HasModified() {
				buf.WriteString(modifiedFmt(modifiedGlyph))
			} else {
				buf.WriteRune(' ')
			}
			// TODO star glyph
			return buf.String()
		}(),
		// dirty/clean
		func() string {
			if pi.IsDirty() {
				return dirtyFmt(dirtyGlyph)
			} else {
				return cleanFmt(cleanGlyph)
			}
		}(),
	)
}

func Run(cwd string) *PorcInfo {
	gitOut, err := GetGitOutput(cwd)
	if err != nil {
		log.Printf("error: %s", err)
		if err == ErrNotAGitRepo {
			os.Exit(0)
		}
		fmt.Printf("error: %s", err)
		os.Exit(1)
	}

	var porcInfo = new(PorcInfo)
	porcInfo.workingDir = cwd

	if err := porcInfo.ParsePorcInfo(gitOut); err != nil {
		log.Printf("error: %s", err)
		fmt.Printf("error: %s", err)
		os.Exit(1)
	}

	return porcInfo
}
