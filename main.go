package main

import (
    "bufio"
    "fmt"
    "os"
    "os/exec"
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
            Git.branch = inf[1]
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
    //fmt.Printf("%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v\n",
    fmt.Printf("commit: %v\nbranch: %v\nremote: %v\nahead: %v\nbehind: %v\nuntr %v\nadd %v\nmod %v\ndel %v\nren %v\ncop %v\n",
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
    stdout, err := cmd.StdoutPipe()

    if err != nil {
        fmt.Fprintln(os.Stderr, "[!]", err)
        return
    }
    if err := cmd.Start(); err != nil {
        fmt.Fprintln(os.Stderr, "[!]", err)
        return
    }

    scanner := bufio.NewScanner(stdout)
    fmt.Println("Opened Pipe. Waiting for git...")

    stop := make(chan bool)
    go readGitStdout(scanner, stop)
    <-stop
    cmd.Wait()

    shellOutput()
    // fmt.Println(Git)
}
