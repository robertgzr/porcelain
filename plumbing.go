// Package plumbing can extract the state of a git repo
// by running various git commands and looking at the .git directory.
//
// That includes:
// - checked out branch
// - HEAD
// - how far ahead/behind the tracked remote branch we are
// - do we have untracked files?
// - do we have unmerged files?
// - do we have unstaged changes?
// - do we have staged but uncommited changes?
//
package porcelain

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"text/scanner"
)

var gitbin string

func Check() (err error) {
	gitbin, err = exec.LookPath("git")
	if err != nil {
		return ErrBinaryNotFound
	}

	c := exec.Command(gitbin, "version")
	buf, err := c.Output()
	if err != nil {
		return wrapErr(ErrRunningGit, err)
	}

	p := NewParser(bytes.NewReader(buf))
	v, err := p.ParseVersion()
	if err != nil {
		return wrapErr(ErrParser, err)
	}

	if !isMinimumVersion(v) {
		return ErrWrongVersion
	}
	return nil
}

func CurrentBranch() (string, error) { return "FAIL", nil }
func CurrentCommit() (string, error) { return "FAIL", nil }
func CommitsAhead() (int, error)     { return 0, nil }
func CommitsBehind() (int, error)    { return 0, nil }
func HasUntracked() (bool, error)    { return false, nil }
func HasUnmerged() (bool, error)     { return false, nil }
func HasUnstaged() (bool, error)     { return false, nil }
func HasStaged() (bool, error)       { return false, nil }

var (
	ErrBinaryNotFound = errors.New("git binary not found")
	ErrWrongVersion   = errors.New("git version <2.13.2")
	ErrRunningGit     = errors.New("error running git")
	ErrParser         = errors.New("failed to parse")
)

type (
	GitVersion = [3]int
)

func wrapErr(outer, inner error) error {
	return fmt.Errorf("%s: %s", outer, inner)
}

func isMinimumVersion(v GitVersion) bool {
	if v[0] < 2 {
		return false
	}
	if v[0] == 2 && v[1] < 13 {
		return false
	}
	if v[0] == 2 && v[1] == 13 && v[2] < 2 {
		return false
	}
	return true
}

// Parser understands git command output.
type Parser interface {
	ParseVersion() (GitVersion, error)
}

// gitParser implements Parser
type gitParser struct {
	rawbuf io.Reader
	scanner.Scanner
	lastError error
}

func NewParser(buf io.Reader) Parser {
	p := gitParser{rawbuf: buf}
	p.Init(buf)
	p.Error = func(_ *scanner.Scanner, msg string) { p.lastError = errors.New(msg) }
	return &p
}

// ParseVersion parser the git command version number and outputs it as an array with 3 uint integers
// Expected input: `git [--version|version]`
//
func (p *gitParser) ParseVersion() (GitVersion, error) {
	var ver = GitVersion([3]int{0, 0, 0})
	var idx = 0

	p.Mode ^= scanner.ScanFloats
	p.Mode |= scanner.ScanInts
	for tok := p.Scan(); tok != scanner.EOF; tok = p.Scan() {
		if tok == scanner.Int {
			i, _ := strconv.Atoi(p.TokenText())
			ver[idx] = i
		}
		if p.TokenText() == "." {
			idx += 1
			continue
		}
	}
	return ver, p.lastError
}
