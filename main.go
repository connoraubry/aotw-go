package main

import (
	"flag"

	"github.com/connoraubry/aotw_go/src/server"
	"github.com/sirupsen/logrus"
)

var (
	dbFile = flag.String("db-path", "./test.db", "Path to the db file")
	// logLevel = flag.String("--log-level", "INFO", "Log level [TRACE, DEBUG, INFO, WARNING]")
	fBool = flag.Bool("debug", false, "Turn debug mode on")
)

func main() {

	flag.Parse()

	if *fBool {
		logrus.SetLevel(logrus.DebugLevel)
	}

	s := server.New(*dbFile)
	s.Run()
}
