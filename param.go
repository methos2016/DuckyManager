package main

// Param holds the parameters for a script.
// They are basically pair values: a name, which is the string to be replaced on the script,
// and a description of said variable
type Param struct {
	Name string
	Desc string
}
