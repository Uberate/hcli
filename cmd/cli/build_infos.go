package main

import "fmt"

var Version = ""
var HashTag = ""
var BranchName = ""
var BuildDate = ""
var GoVersion = ""

func ShowVersion() {
	fmt.Printf(`Version: %s
HashTag: %s
BranchName: %s
BuildDate: %s
GoVersion: %s`, Version, HashTag, BranchName, BuildDate, GoVersion)
}
