package main

import "log"

var (
	translate Strings
	config    Config
	l         *log.Logger
)

const (
	languageVer = "0.3"
	languageDir = "language"
	configFile  = "config.json"

	errStr = "[-] "
	okStr  = "[+] "

	errExitCode = -1
)
