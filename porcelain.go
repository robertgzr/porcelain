package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/robertgzr/color"
)

// TODO allow custom log location
const logloc string = "/tmp/porcelain.log"

var (
	cwd                string
	debugFlag, fmtFlag bool
)

type GitArea struct {
	modified int
	added    int
	deleted  int
	renamed  int
	copied   int
}

func (a *GitArea) hasChanged() bool {
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

func (pi *PorcInfo) hasUnmerged() bool {
	if pi.unmerged > 0 {
		return true
	}
	gitDir, err := PathToGitDir(cwd)
	if err != nil {
		log.Println(cwd, err)
		return false
	}
	// TODO figure out if output of MERGE_HEAD can be useful
	if _, err := ioutil.ReadFile(path.Join(gitDir, "MERGE_HEAD")); err != nil {
		if os.IsNotExist(err) {
			return false
		}
		log.Println(err)
		return false
	} else {
		return true
	}
}
func (pi *PorcInfo) hasModified() bool {
	return pi.Unstaged.hasChanged()
}
func (pi *PorcInfo) isDirty() bool {
	return pi.Staged.hasChanged()
}

func (pi *PorcInfo) Debug() string {
	return fmt.Sprintf("%#+v", pi)
}

// Fmt formats the output for the shell
// TODO should be configurable by the user
//
func (pi *PorcInfo) Fmt() string {
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

	color.NoColor = false
	color.EscapeZshPrompt = true

	branchFmt := color.New(color.FgBlue).SprintFunc()
	commitFmt := color.New(color.FgGreen, color.Italic).SprintFunc()

	aheadFmt := color.New(color.Faint, color.BgCyan, color.FgBlack).SprintFunc()
	behindFmt := color.New(color.Faint, color.BgRed, color.FgWhite).SprintFunc()

	modifiedFmt := color.New(color.FgBlue).SprintFunc()
	// deletedFmt := color.New(color.FgYellow).SprintFunc()
	dirtyFmt := color.New(color.FgRed).SprintFunc()
	cleanFmt := color.New(color.FgGreen).SprintFunc()

	untrackedFmt := color.New(color.Faint).SprintFunc()
	unmergedFmt := color.New(color.FgYellow).SprintFunc()

	return fmt.Sprintf("%s %s@%s %s %s %s",
		branchGlyph,
		branchFmt(pi.branch),
		commitFmt(pi.commit[:7]),
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
			if pi.hasUnmerged() {
				buf.WriteString(unmergedFmt(unmergedGlyph))
			} else {
				buf.WriteRune(' ')
			}
			if pi.hasModified() {
				buf.WriteString(modifiedFmt(modifiedGlyph))
			} else {
				buf.WriteRune(' ')
			}
			// TODO star glyph
			return buf.String()
		}(),
		// dirty/clean
		func() string {
			if pi.isDirty() {
				return dirtyFmt(dirtyGlyph)
			} else {
				return cleanFmt(cleanGlyph)
			}
		}(),
	)
}

func init() {
	flag.BoolVar(&debugFlag, "debug", false, "print output for debugging")
	flag.BoolVar(&fmtFlag, "fmt", false, "print formatted output")
	flag.Parse()

	logFd, err := os.OpenFile(logloc, os.O_CREATE|os.O_RDWR|os.O_APPEND, os.ModePerm)
	if err != nil {
		os.Exit(1)
	}
	log.SetOutput(logFd)
	log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)

	cwd, _ = os.Getwd()
}

func run() *PorcInfo {
	gitOut, err := GetGitOutput(cwd)
	if err != nil {
		if err == ErrNotAGitRepo {
			os.Exit(0)
		}
		fmt.Print("sry, no info :(")
		os.Exit(1)
	}

	var porcInfo = new(PorcInfo)
	if err := porcInfo.ParsePorcInfo(gitOut); err != nil {
		fmt.Print("sry, no info :(")
		os.Exit(1)
	}

	return porcInfo
}

func main() {
	var out string
	switch {
	case debugFlag:
		out = run().Debug()
	case fmtFlag:
		out = run().Fmt()
	default:
		flag.Usage()
		fmt.Println("\nOutside of a repository there will be no output.")
		os.Exit(1)
	}

	fmt.Fprint(os.Stdout, out)
}
