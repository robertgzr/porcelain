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
	"os/exec"
	"strconv"
	"strings"
)

var (
	gitbin       string
	gitdir       string
	gitporcelain []byte
)

var (
	ErrBinaryNotFound = errors.New("git binary not found")
	ErrWrongVersion   = errors.New("git version <2.13.2")
	ErrRunningGit     = errors.New("error running git")
	ErrUnchecked      = errors.New("attempted to do unchecked plumbing, call CheckGit first")
)

func CheckGit() (err error) {
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
	version, err := p.ParseVersion()
	if err != nil {
		return wrapErr(ErrParser, err)
	}

	if !isMinimumVersion(version) {
		return ErrWrongVersion
	}
	return nil
}

func CheckDir(path string) error {
	c := exec.Command(gitbin, "rev-parse", "--is-inside-work-tree")
	c.Dir = path
	buf, err := c.Output()
	if err != nil {
		return err
	}
	if ok, err := strconv.ParseBool(strings.TrimSpace(string(buf))); err != nil || !ok {
		return err
	}
	gitdir = path
	return nil
}

func CurrentBranch() (string, error) {
	if gitporcelain == nil {
		if err := getPorcelain(); err != nil {
			return "", err
		}
	}

	p := NewParser(bytes.NewReader(gitporcelain))
	return p.ParseBranch()
}
func CurrentCommit() (string, error) {
	if gitporcelain == nil {
		if err := getPorcelain(); err != nil {
			return "", err
		}
	}

	p := NewParser(bytes.NewReader(gitporcelain))
	return p.ParseCommit()
}
func CommitsAhead() (int, error) {
	if gitporcelain == nil {
		if err := getPorcelain(); err != nil {
			return 0, err
		}
	}

	p := NewParser(bytes.NewReader(gitporcelain))
	return p.ParseAhead()
}
func CommitsBehind() (int, error) {
	if gitporcelain == nil {
		if err := getPorcelain(); err != nil {
			return 0, err
		}
	}

	p := NewParser(bytes.NewReader(gitporcelain))
	return p.ParseBehind()
}
func HasUntracked() (bool, error) { return false, nil }
func HasUnmerged() (bool, error)  { return false, nil }
func HasUnstaged() (bool, error)  { return false, nil }
func HasStaged() (bool, error)    { return false, nil }

func getPorcelain() (err error) {
	if gitdir == "" {
		return ErrUnchecked
	}
	c := exec.Command(gitbin, "status", "--porcelain=v2", "--branch")
	c.Dir = gitdir
	gitporcelain, err = c.Output()
	return
}

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
