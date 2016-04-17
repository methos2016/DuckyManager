package main

import "log"

var (
	translate Strings
	config    Config
	l         *log.Logger

	debug bool

	lang string
)

const (
	languageVer = "0.4"
	languageDir = "language"
	configFile  = "config.json"

	errStr  = "[-] "
	okStr   = "[+] "
	infoStr = "[i] "

	errExitCode = -1
)
