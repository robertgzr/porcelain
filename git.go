package main

import (
	"bytes"
	"errors"
	"io"
	"os/exec"
	"strings"
)

const notRepoStatus string = "exit status 128"

var ErrNotAGitRepo error = errors.New("not a git repo")

func GetGitOutput(cwd string) (io.Reader, error) {
	var buf = new(bytes.Buffer)

	cmd := exec.Command("git", "status", "--porcelain=v2", "--branch")
	cmd.Stderr = buf
	cmd.Stdout = buf

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
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}
