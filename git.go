package main

import (
	"bytes"
	"errors"
	"io"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

const notRepoStatus string = "exit status 128"

var ErrNotAGitRepo error = errors.New("not a git repo")

func GetGitOutput(cwd string) (io.Reader, error) {
	if ok, err := IsInsideWorkTree(cwd); err != nil {
		log.Printf("error detecting work tree: %s", err)
		return nil, err
	} else if !ok {
		return nil, ErrNotAGitRepo
	}

	var buf = new(bytes.Buffer)
	cmd := exec.Command("git", "status", "--porcelain=v2", "--branch")
	cmd.Stdout = buf
	cmd.Dir = cwd
	log.Printf("running %q", cmd.Args)

	if err := cmd.Run(); err != nil {
		if cmd.ProcessState.String() == notRepoStatus {
			return nil, ErrNotAGitRepo
		}
		return nil, err
	}

	return buf, nil
}

func PathToGitDir(cwd string) (string, error) {
	cmd := exec.Command("git", "rev-parse", "--absolute-git-dir")
	cmd.Dir = cwd
	log.Printf("running %q", cmd.Args)

	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func IsInsideWorkTree(cwd string) (bool, error) {
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	cmd.Dir = cwd
	log.Printf("running %q", cmd.Args)

	out, err := cmd.Output()
	if err != nil {
		return false, err
	}
	return strconv.ParseBool(strings.TrimSpace(string(out)))
}
