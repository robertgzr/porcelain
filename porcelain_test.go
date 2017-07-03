package main

import (
	"strings"
	"testing"
)

const expectedDebugOutput string = "&main.PorcInfo{branch:\"master\", commit:\"51c9c58e2175b768137c1e38865f394c76a7d49d\", remote:\"\", upstream:\"origin/master\", ahead:1, behind:10, untracked:5, unmerged:1, Unstaged:main.GitArea{modified:3, added:0, deleted:1, renamed:0, copied:0}, Staged:main.GitArea{modified:0, added:0, deleted:0, renamed:1, copied:0}}"

func TestDebugOutput(t *testing.T) {
	var pi = new(PorcInfo)
	if err := pi.ParsePorcInfo(strings.NewReader(gitoutput)); err != nil {
		t.Fatal(err)
	}

	if out := pi.Debug(); out != expectedDebugOutput {
		t.Logf("\nexpected:\n%s\ngot:\n%s\n", expectedDebugOutput, out)
		t.FailNow()
	}
}
