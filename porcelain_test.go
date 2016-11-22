package main

import (
	"fmt"
	"reflect"
	"testing"
)

var (
	branchInfoStrings = []string{
		"## new2",
		"## master...origin/master",
		"## 0.9...test/0.9",
		"## master...origin/master [ahead 1]",
		"## master...origin/master [ahead 1, behind 99]",
		"## Initial commit on master",
		"## HEAD (no branch)",
	}
	branchInfoExpected = []gitinfo{
		gitinfo{
			branch: "new2", commit: "", remote: "", trackedBranch: "", ahead: 0, behind: 0},
		gitinfo{
			branch: "master", commit: "", remote: "origin", trackedBranch: "origin/master", ahead: 0, behind: 0},
		gitinfo{
			branch: "0.9", commit: "", remote: "test", trackedBranch: "test/0.9", ahead: 0, behind: 0},
		gitinfo{
			branch: "master", commit: "", remote: "origin", trackedBranch: "origin/master", ahead: 1, behind: 0},
		gitinfo{
			branch: "master", commit: "", remote: "origin", trackedBranch: "origin/master", ahead: 1, behind: 99},
		gitinfo{
			branch: "master", commit: "", remote: "", trackedBranch: "", ahead: 0, behind: 0},
		gitinfo{
			branch: "", commit: "HEAD", remote: "", trackedBranch: "", ahead: 0, behind: 0},
	}
)

func TestStatusParser(t *testing.T) {
	for i, s := range branchInfoStrings {
		t.Run(fmt.Sprintf("TestStatusParser_%d", i), func(t *testing.T) {
			parseBranchinfo(s)
			t.Logf("Parsed... '%s'", s)
			if !reflect.DeepEqual(branchInfoExpected[i], Git) {
				t.Fatalf("\n\texp: %+v\n\tgot: %+v", branchInfoExpected[i], Git)
			}
			Git = gitinfo{}
		})
	}
}
