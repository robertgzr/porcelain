package main // import "github.com/robertgzr/porcelain/cmd/porcelain"

import (
	"flag"
	"io/ioutil"
	"log"
	"os"

	"github.com/robertgzr/porcelain"
	"github.com/robertgzr/porcelain/fmt"
)

const (
	defaultDebug        = false
	defaultNoColor      = false
	defaultEscapeBash   = false
	defaultEscapeZsh    = false
	defaultEscapeTmux   = false
	defaultPrintVersion = false
	defaultLogToStderr  = false
	defaultLogToFile    = "/tmp/porcelain.log"
	defaultPath         = "."
)

var (
	flagDebug        bool
	flagNoColor      bool
	flagEscapeBash   bool
	flagEscapeZsh    bool
	flagEscapeTmux   bool
	flagPrintVersion bool
	flagLogToStderr  bool
	flagLogToFile    string
	flagPath         string
)

func main() {
	flag.BoolVar(&flagDebug, "debug", defaultDebug, "print logs")
	flag.BoolVar(&flagNoColor, "no-color", defaultNoColor, "print formatted output without color codes")

	flag.BoolVar(&flagEscapeBash, "bash", defaultEscapeBash, "escape fmt output for bash")
	flag.BoolVar(&flagEscapeZsh, "zsh", defaultEscapeZsh, "escape fmt output for zsh")
	flag.BoolVar(&flagEscapeTmux, "tmux", defaultEscapeTmux, "escape fmt output for tmux")

	flag.BoolVar(&flagPrintVersion, "version", defaultPrintVersion, "print version and exit")
	flag.BoolVar(&flagLogToStderr, "logtostderr", defaultLogToStderr, "write logs to stderr")
	flag.StringVar(&flagLogToFile, "logtofile", defaultLogToFile, "write logs to a file")

	flag.StringVar(&flagPath, "path", defaultPath, "show output for path instead of the working directory")

	flag.Bool("fmt", true, "print formatted output (compat only - does nothing)")
	flag.Parse()

	// set up output and logging
	stdout := log.New(os.Stdout, "", 0)
	log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)
	log.SetOutput(ioutil.Discard)

	if flagPrintVersion {
		stdout.Printf("porcelain version %s (%s)\nbuilt %s\n",
			porcelain.BuildVersion, porcelain.BuildCommit, porcelain.BuildTimestamp)
		os.Exit(0)
	}

	if flagDebug {
		logFile, err := os.OpenFile(flagLogToFile, os.O_CREATE|os.O_RDWR|os.O_APPEND, os.ModePerm)
		if err != nil {
			stdout.Fatalf("failed to open log file at %s", flagLogToFile)
		}
		log.SetOutput(logFile)
	}

	if flagLogToStderr {
		log.SetOutput(os.Stderr)
	}

	if flagPath == "." {
		var err error
		flagPath, err = os.Getwd()
		if err != nil {
			log.Fatal("failed to get current working directory, please use --path")
		}
	}

	if err := porcelain.CheckGit(); err != nil {
		log.Fatal(err)
	}

	if err := porcelain.CheckDir(flagPath); err != nil {
		log.Fatal(err)
	}

	var (
		err error
		rs  porcelain.RepoState
	)

	rs.Branch, err = porcelain.CurrentBranch()
	if err != nil {
		log.Fatal(err)
	}
	rs.Commit, err = porcelain.CurrentCommit()
	if err != nil {
		log.Fatal(err)
	}
	rs.CommitsAhead, err = porcelain.CommitsAhead()
	if err != nil {
		log.Fatal(err)
	}
	rs.CommitsBehind, err = porcelain.CommitsBehind()
	if err != nil {
		log.Fatal(err)
	}
	rs.Untracked, err = porcelain.HasUntracked()
	if err != nil {
		log.Fatal(err)
	}
	rs.Unmerged, err = porcelain.HasUnmerged()
	if err != nil {
		log.Fatal(err)
	}
	rs.Unstaged, err = porcelain.HasUnstaged()
	if err != nil {
		log.Fatal(err)
	}
	rs.Staged, err = porcelain.HasStaged()
	if err != nil {
		log.Fatal(err)
	}

	formatter := fmt.New()
	formatter.NoColor = flagNoColor
	formatter.ModeBash = flagEscapeBash
	formatter.ModeZsh = flagEscapeZsh
	formatter.ModeTmux = flagEscapeTmux
	formatter.Format(os.Stdout, rs)
}
