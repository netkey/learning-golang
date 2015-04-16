package main

import (
	"flag"
	"fmt"
	"github.com/ezotrank/logsend/logsend"
	logpkg "log"
	"os"
)

const (
	VERSION = "1.7.1"
)

var (
	watchDir      = flag.String("watch-dir", "", "deprecated, simply add the directory as an argument, in the end")
	config        = flag.String("config", "", "path to config.json file")
	check         = flag.Bool("check", false, "check config.json")
	debug         = flag.Bool("debug", false, "turn on debug messages")
	continueWatch = flag.Bool("continue-watch", false, "watching folder for new files")
	logFile       = flag.String("log", "", "log file")
	dryRun        = flag.Bool("dry-run", false, "not send data")
	memprofile    = flag.String("memprofile", "", "memory profiler")
	maxprocs      = flag.Int("maxprocs", 0, "max count of cpu")
	readWholeLog  = flag.Bool("read-whole-log", false, "read whole logs")
	readOnce      = flag.Bool("read-once", false, "read logs once and exit")
	regex         = flag.String("regex", "", "regex rule")
	version       = flag.Bool("version", false, "show version number")
)

func main() {
	flag.Parse()

	if *version {
		fmt.Printf("logsend version %v\n", VERSION)
		os.Exit(0)
	}

	if *logFile != "" {
		file, err := os.OpenFile(*logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			fmt.Errorf("Failed to open log file: %+v\n", err)
		}
		defer file.Close()
		logsend.Conf.Logger = logpkg.New(file, "", logpkg.Ldate|logpkg.Ltime|logpkg.Lshortfile)
	}

	logsend.Conf.Debug = *debug
	logsend.Conf.ContinueWatch = *continueWatch
	logsend.Conf.WatchDir = *watchDir
	logsend.Conf.Memprofile = *memprofile
	logsend.Conf.DryRun = *dryRun
	logsend.Conf.ReadWholeLog = *readWholeLog
	logsend.Conf.ReadOnce = *readOnce

	if *check {
		_, err := logsend.LoadConfigFromFile(*config)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("ok")
		os.Exit(0)
	}

	fi, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}

	logDirs := make([]string, 0)
	if len(flag.Args()) > 0 {
		logDirs = flag.Args()
	} else {
		logDirs = append(logDirs, *watchDir)
	}

	if fi.Mode()&os.ModeNamedPipe == 0 {
		logsend.WatchFiles(logDirs, *config)
	} else {
		flag.VisitAll(logsend.LoadRawConfig)
		logsend.ProcessStdin()
	}
	os.Exit(0)
}
