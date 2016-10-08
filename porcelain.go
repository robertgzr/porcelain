package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

const gitbin string = "/usr/bin/git"

var (
	debugFlag, basicFlag, fmtFlag bool

	gitrevparse = []string{"rev-parse", "--short", "HEAD"}
	gitstatus   = []string{"status", "--porcelain", "--branch"}
)

type gitinfo struct {
	branch        string
	commit        string
	remote        string
	trackedBranch string
	ahead         int
	behind        int

	untracked int // ?
	dirty     int // changes not in index

	modified int
	added    int
	deleted  int
	renamed  int
	copied   int
	unmerged int // diff flag
}

var Git gitinfo

func sliceContains(sl []string, cmp string) int {
	for i, a := range sl {
		if a == cmp {
			return i
		}
	}
	return -1
}

func parseBranchinfo(s string) {
	var (
		matchBranch  []string
		matchDiffers []string
		err          error
	)

	reBranchOrigin := regexp.MustCompile("\\s([a-zA-Z0-9-_\\.]+)(?:\\.\\.\\.)([a-zA-Z0-9-_\\.]+)\\/([a-zA-Z0-9-_\\.]+)(.*)|([a-zA-Z0-9-_\\.]+)$")
	matchBranch = reBranchOrigin.FindStringSubmatch(s)

	if matchBranch[2] != "" {
		Git.branch = matchBranch[1]
		Git.remote = matchBranch[2]
		Git.trackedBranch = matchBranch[2] + "/" + matchBranch[3]
	} else {
		Git.branch = matchBranch[5]
		Git.remote = "-"
		Git.trackedBranch = "-"
	}

	// match ahead/behind part
	reDiffers := regexp.MustCompile("[0-9]+")
	matchDiffers = reDiffers.FindAllString(matchBranch[4], 2)

	switch len(matchDiffers) {
	case 2:
		Git.behind, err = strconv.Atoi(matchDiffers[1])
		if err != nil {
			panic(err)
		}
		fallthrough
	case 1:
		Git.ahead, err = strconv.Atoi(matchDiffers[0])
		if err != nil {
			panic(err)
		}
	default:
		Git.behind = 0
		Git.ahead = 0
	}
}

func parseLine(line string) {
	switch line[:2] {

	// match branch and origin
	case "##":
		parseBranchinfo(line)

	// untracked files
	case "??":
		Git.untracked++

	case "MM":
		fallthrough
	case "AM":
		fallthrough
	case "RM":
		fallthrough
	case "CM":
		fallthrough
	case " M":
		Git.modified++
		Git.dirty++

	case "MD":
		fallthrough
	case "AD":
		fallthrough
	case "RD":
		fallthrough
	case "CD":
		fallthrough
	case " D":
		Git.deleted++
		Git.dirty++

	// changes in the index
	case "M ":
		Git.modified++
	case "A ":
		Git.added++
	case "D ":
		Git.deleted++
	case "R ":
		Git.renamed++
	case "C ":
		Git.copied++

	case "DD":
		fallthrough
	case "AU":
		fallthrough
	case "UD":
		fallthrough
	case "UA":
		fallthrough
	case "DU":
		fallthrough
	case "AA":
		fallthrough
	case "UU":
		Git.unmerged++

	// catch everything else
	default:
		fmt.Println(line)
		panic("unexpected input.")
	}
}

func readGitStdout(scanner *bufio.Scanner, stop chan bool) {
	for scanner.Scan() {
		line := scanner.Text()
		parseLine(line)
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "[!]", err)
	}
	stop <- true
}

func basicOutput() {
	fmt.Printf("%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v\n",
		Git.commit,
		Git.branch,
		Git.trackedBranch,
		Git.ahead,
		Git.behind,
		Git.untracked,
		Git.added,
		Git.modified,
		Git.deleted,
		Git.renamed,
		Git.copied)
}

func debugOutput() {
	fmt.Printf("%+v\n", Git)
}

func formattedOutput() {
	var (
		branchGlyph    string = ""
		modifiedGlyph  string = "Δ"
		deletedGlyph   string = "＊"
		dirtyGlyph     string = "✘"
		cleanGlyph     string = "✔"
		untrackedGlyph string = "?"
		unmergedGlyph  string = "‼"
		aheadArrow     string = "↑"
		behindArrow    string = "↓"
	)

	branchFmt := color.New(color.FgHiBlue).SprintFunc()
	commitFmt := color.New(color.FgHiGreen, color.Italic).SprintFunc()

	aheadFmt := color.New(color.Faint, color.BgCyan, color.FgBlack).SprintFunc()
	behindFmt := color.New(color.Faint, color.BgHiRed, color.FgWhite).SprintFunc()

	modifiedFmt := color.New(color.FgBlue).SprintFunc()
	deletedFmt := color.New(color.FgYellow).SprintFunc()
	dirtyFmt := color.New(color.FgHiRed).SprintFunc()
	cleanFmt := color.New(color.FgGreen).SprintFunc()

	untrackedFmt := color.New(color.Faint).SprintFunc()
	unmergedFmt := color.New(color.BgMagenta, color.FgHiWhite).SprintFunc()

	fmt.Printf("%s %s@%s %s%s %s%s %s%s %s",
		branchGlyph,
		branchFmt(Git.branch),
		commitFmt(Git.commit),
		//ahead/behind
		func(n int) string {
			if n > 0 {
				return aheadFmt(" ", aheadArrow, n, " ")
			} else {
				return ""
			}
		}(Git.ahead),
		func(n int) string {
			if n > 0 {
				return behindFmt(" ", behindArrow, n, " ")
			} else {
				return ""
			}
		}(Git.behind),
		// stats
		// untracked
		func(n int) string {
			if n > 0 {
				return untrackedFmt(untrackedGlyph)
			} else {
				return ""
			}
		}(Git.untracked),
		// unmerged
		func(n int) string {
			if n > 0 {
				return unmergedFmt(unmergedGlyph)
			} else {
				return ""
			}
		}(Git.unmerged),
		// modi
		func(n int) string {
			if n > 0 {
				return modifiedFmt(modifiedGlyph)
			} else {
				return ""
			}
		}(Git.modified),
		// del
		func(n int) string {
			if n > 0 {
				return deletedFmt(deletedGlyph)
			} else {
				return ""
			}
		}(Git.deleted),
		// dirty/clean
		func(n int) string {
			if n > 0 {
				return dirtyFmt(dirtyGlyph)
			} else {
				return cleanFmt(cleanGlyph)
			}
		}(Git.dirty),
	)
}

func execRevParse() string {
	// commit
	cmd := exec.Command(gitbin, gitrevparse...)
	out, err := cmd.Output()
	if err != nil {
		if strings.Contains(err.Error(), "128") {
			return "initial"
		} else {
			panic(err)
		}
	}

	return string(out)
}

func execStatus() {
	cmd := exec.Command(gitbin, gitstatus...)
	stdout, err := cmd.StdoutPipe()
	// catch pipe errors
	if err != nil {
		panic(err)
	}
	if err := cmd.Start(); err != nil {
		panic(err)
	}

	stop := make(chan bool)
	go readGitStdout(bufio.NewScanner(stdout), stop)
	<-stop
	cmd.Wait()

}

func init() {
	flag.BoolVar(&debugFlag, "debug", false, "print output for debugging")
	flag.BoolVar(&basicFlag, "basic", false, "print basic number output")
	flag.BoolVar(&fmtFlag, "fmt", false, "print formatted output")
	flag.Parse()
}

func main() {
	out := execRevParse()
	Git.commit = strings.TrimSuffix(string(out), "\n")

	execStatus()

	switch {
	case debugFlag:
		debugOutput()
	case basicFlag:
		basicOutput()
	case fmtFlag:
		formattedOutput()
	default:
		flag.Usage()
	}
}
