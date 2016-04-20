package main

import (
	"log"
	"time"
)

var (
	translate Strings
	config    Config
	l         *log.Logger

	debug bool

	lang string
)

const (
	languageVer = "0.6"
	languageDir = "language"
	configFile  = "config.json"

	errStr = "[-] "
	okStr  = "[+] "
	//infoStr = "[i] "

	githubRFC = time.RFC3339 // Close RFC as the defined by github's API

	errExitCode = -1
)
