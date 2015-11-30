package main

import (
    "bufio"
    "fmt"
    "os"
    "os/exec"
    "regexp"
    "strings"
)

var Git struct {
    branch    string
    commit    string
    remote    string
    ahead     int
    behind    int
    untracked int // ?
    added     int // A
    modified  int // M
    deleted   int // D
    renamed   int // R
    copied    int // C
    unmerged  int // U
}

func parseLine(line string) {
    inf := strings.Fields(line)
    if strings.Contains(inf[0], "#") {
        // branch info
        if strings.Contains(line, "Initial") {
            Git.branch = "master"
            Git.commit = "init"
        } else {
            re := regexp.MustCompile("([a-zA-Z0-9]+)").FindAllString(inf[1], -1)
            Git.branch = re[0]
            Git.remote = re[1]
            // todo:
            // parse remote, ahead, behind
        }
    }
    if strings.Contains(inf[0], "?") {
        // untracked files
        Git.untracked++
    }
    if strings.Contains(inf[0], "A") {
        // added files
        Git.added++
    }
    if strings.Contains(inf[0], "M") {
        // modified files
        Git.modified++
    }
    if strings.Contains(inf[0], "D") {
        // deleted files
        Git.deleted++
    }
    if strings.Contains(inf[0], "R") {
        // renamed files
        Git.renamed++
    }
    if strings.Contains(inf[0], "C") {
        // copied files
        Git.copied++
    }
    if strings.Contains(inf[0], "U") {
        // unmerged files
        Git.unmerged++
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

func shellOutput() {
    fmt.Printf("%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v\n",
        //fmt.Printf("commit: %v\nbranch: %v\nremote: %v\nahead: %v\nbehind: %v\nuntr %v\nadd %v\nmod %v\ndel %v\nren %v\ncop %v\n",
        Git.commit,
        Git.branch,
        Git.remote,
        Git.ahead,
        Git.behind,
        Git.untracked,
        Git.added,
        Git.modified,
        Git.deleted,
        Git.renamed,
        Git.copied)
}

func main() {
    cmd := exec.Command("/usr/local/bin/git", "status", "--porcelain", "--branch")
    cmd2 := exec.Command("/usr/local/bin/git", "rev-parse", "--short", "HEAD")

    stdout, err := cmd.StdoutPipe()
    // catch pipe errors
    if err != nil {
        fmt.Fprintln(os.Stderr, "[!]", err)
        return
    }

    // fork child
    // catch fork errors
    if err := cmd.Start(); err != nil {
        fmt.Fprintln(os.Stderr, "[!]", err)
        return
    }
    // commit
    out, err := cmd2.Output()
    if err != nil {
        fmt.Fprintln(os.Stderr, "[!]", err)
        return
    }
    Git.commit = strings.TrimSuffix(string(out), "\n")

    stop := make(chan bool)
    go readGitStdout(bufio.NewScanner(stdout), stop)
    <-stop
    cmd.Wait()

    shellOutput()
    // fmt.Println(Git)
}
