package main

import (
    "bufio"
    "flag"
    "fmt"
    "os"
    "os/exec"
    "regexp"
    "strconv"
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

func SliceContains(sl []string, cmp string) int {
    for i, a := range sl {
        if a == cmp {
            return i
        }
    }
    return -1
}

func parseLine(line string) {
    inf := strings.Fields(line)
    if strings.Contains(inf[0], "#") {
        // branch info
        if strings.Contains(line, "Initial") {
            Git.branch = "master"
            Git.commit = "init"
        } else {
            re := regexp.MustCompile("([a-zA-Z0-9-_]+)").FindAllString(line, -1)
            if re[0] != "" {
                Git.branch = re[0]
            }
            if re[1] != "" {
                Git.remote = re[1]
            }
            if i := SliceContains(re, "ahead"); i != -1 {
                Git.ahead, _ = strconv.Atoi(re[i+1])
            }
            if i := SliceContains(re, "behind"); i != -1 {
                Git.behind, _ = strconv.Atoi(re[i+1])
            }
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

func debugOutput() {
    fmt.Printf("commit:\t%v\nbranch:\t%v\nremote:\t%v\nahead:\t%v\nbehind:\t%v\nuntr:\t%v\nadd:\t%v\nmod:\t%v\ndel:\t%v\nren:\t%v\ncop:\t%v\n",
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
    debug := flag.Bool("debug", false, "print output for debugging")
    flag.Parse()

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

    // print debug output if -debug flag is set
    if *debug == false {
        shellOutput()
    } else {
        fmt.Printf("go-gitparser v1.1 Debug mode:\n\n%v\n", Git)
        debugOutput()
    }
}
