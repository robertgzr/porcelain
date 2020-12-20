package porcelain

import (
	"bufio"
	"io"
	"log"
	"strconv"
	"strings"
)

func consumeNext(s *bufio.Scanner) string {
	if s.Scan() {
		return s.Text()
	}
	return ""
}

func (pi *PorcInfo) ParsePorcInfo(r io.Reader) error {
	log.Println("parsing git output")

	var err error
	var s = bufio.NewScanner(r)

	for s.Scan() {
		if len(s.Text()) < 1 {
			continue
		}

		// NOTE we only return the last error, probably should fail early?
		err = pi.ParseLine(s.Text())
		if err != nil {
			log.Printf("error parsing %q: %v", s.Text(), err)
		}
	}

	return err
}

func (pi *PorcInfo) ParseLine(line string) error {
	s := bufio.NewScanner(strings.NewReader(line))
	// switch to a word based scanner
	s.Split(bufio.ScanWords)

	var err error
	for s.Scan() {
		switch s.Text() {
		case "#":
			err = pi.parseBranchInfo(s)
			if err != nil {
				log.Printf("error parsing branch: %v", err)
			}
		case "1":
			err = pi.parseTrackedFile(s)
			if err != nil {
				log.Printf("error parsing tracked file: %v", err)
			}
		case "2":
			err = pi.parseRenamedFile(s)
			if err != nil {
				log.Printf("error parsing renamed file: %v", err)
			}
		case "u":
			pi.unmerged++
		case "?":
			pi.untracked++
		}
	}
	return nil
}

func (pi *PorcInfo) parseBranchInfo(s *bufio.Scanner) (err error) {
	// uses the word based scanner from ParseLine
	for s.Scan() {
		switch s.Text() {
		case "branch.oid":
			pi.commit = consumeNext(s)
		case "branch.head":
			pi.branch = consumeNext(s)
		case "branch.upstream":
			pi.upstream = consumeNext(s)
		case "branch.ab":
			err = pi.parseAheadBehind(s)
			if err != nil {
				log.Printf("error parsing branch.ab: %v", err)
			}
		}
	}
	return err
}

func (pi *PorcInfo) parseAheadBehind(s *bufio.Scanner) error {
	// uses the word based scanner from ParseLine
	for s.Scan() {
		i, err := strconv.Atoi(s.Text()[1:])
		if err != nil {
			return err
		}

		switch s.Text()[:1] {
		case "+":
			pi.ahead = i
		case "-":
			pi.behind = i
		}
	}
	return nil
}

// parseTrackedFile parses the porcelain v2 output for tracked entries
// doc: https://git-scm.com/docs/git-status#_changed_tracked_entries
//
func (pi *PorcInfo) parseTrackedFile(s *bufio.Scanner) error {
	// uses the word based scanner from ParseLine
	var index int
	var err error
	for s.Scan() {
		switch index {
		case 0: // xy
			err = pi.parseXY(s.Text())
			if err != nil {
				log.Printf("error parsing XY: %q: %v", s.Text(), err)
			}
		default:
			continue
			// case 1: // sub
			// 	if s.Text() != "N..." {
			// 		log.Println("is submodule!!!")
			// 	}
			// case 2: // mH - octal file mode in HEAD
			// 	log.Println(index, s.Text())
			// case 3: // mI - octal file mode in index
			// 	log.Println(index, s.Text())
			// case 4: // mW - octal file mode in worktree
			// 	log.Println(index, s.Text())
			// case 5: // hH - object name in HEAD
			// 	log.Println(index, s.Text())
			// case 6: // hI - object name in index
			// 	log.Println(index, s.Text())
			// case 7: // path
			// 	log.Println(index, s.Text())
		}
		index++
	}
	return nil
}

func (pi *PorcInfo) parseXY(xy string) error {
	switch xy[:1] { // parse staged
	case "M":
		pi.Staged.modified++
	case "A":
		pi.Staged.added++
	case "D":
		pi.Staged.deleted++
	case "R":
		pi.Staged.renamed++
	case "C":
		pi.Staged.copied++
	}

	switch xy[1:] { // parse unstaged
	case "M":
		pi.Unstaged.modified++
	case "A":
		pi.Unstaged.added++
	case "D":
		pi.Unstaged.deleted++
	case "R":
		pi.Unstaged.renamed++
	case "C":
		pi.Unstaged.copied++
	}
	return nil
}

func (pi *PorcInfo) parseRenamedFile(s *bufio.Scanner) error {
	return pi.parseTrackedFile(s)
}
