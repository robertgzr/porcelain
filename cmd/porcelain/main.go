package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/robertgzr/porcelain"
)

// TODO allow custom log location
const logloc string = "/tmp/porcelain.log"

var (
	commit  string = "invalid"
	version string = "invalid"
	date    string = "invalid"
)

var (
	cwd                     string
	noColorFlag             bool
	debugFlag, fmtFlag      bool
	bashFmtFlag, zshFmtFlag bool
	tmuxFmtFlag             bool
	versionFlag             bool
)

func main() {
	flag.BoolVar(&debugFlag, "debug", false, "write logs to file ("+logloc+")")
	flag.BoolVar(&fmtFlag, "fmt", true, "print formatted output (default)")
	flag.BoolVar(&bashFmtFlag, "bash", false, "escape fmt output for bash")
	flag.BoolVar(&noColorFlag, "no-color", false, "print formatted output without color codes")
	flag.BoolVar(&zshFmtFlag, "zsh", false, "escape fmt output for zsh")
	flag.BoolVar(&tmuxFmtFlag, "tmux", false, "escape fmt output for tmux")
	flag.StringVar(&cwd, "path", "", "show output for path instead of the working directory")
	flag.BoolVar(&versionFlag, "version", false, "print version and exit")

	logtostderr := flag.Bool("logtostderr", false, "write logs to stderr")
	flag.Parse()

	if versionFlag {
		fmt.Printf("porcelain version %s (%s)\nbuilt %s\n", version, commit, date)
		os.Exit(0)
	}

	if debugFlag {
		var (
			err   error
			logFd io.Writer
		)
		if *logtostderr {
			logFd = os.Stderr
		} else {
			logFd, err = os.OpenFile(logloc, os.O_CREATE|os.O_RDWR|os.O_APPEND, os.ModePerm)
			if err != nil {
				os.Exit(1)
			}
		}
		log.SetOutput(logFd)
		log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)
	} else {
		log.SetOutput(ioutil.Discard)
	}

	if cwd == "" {
		cwd, _ = os.Getwd()
	}

	log.Println("running porcelain...")
	log.Println("in directory:", cwd)

	var out string
	switch {
	case fmtFlag:
		out = porcelain.Run(cwd).Fmt(cwd, noColorFlag, bashFmtFlag, zshFmtFlag, tmuxFmtFlag)
	default:
		flag.Usage()
		fmt.Println("\nOutside of a repository there will be no output.")
		os.Exit(1)
	}

	fmt.Fprint(os.Stdout, out)
}
