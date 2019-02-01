package porcelain

type RepoState struct {
	Branch string
	Commit string

	CommitsAhead  int
	CommitsBehind int

	Untracked bool
	Unmerged  bool
	Unstaged  bool
	Staged    bool

	Clean bool
}
