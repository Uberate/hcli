package main

import (
	"context"
	"github.io/uberate/hcli/pkg/outputer"
)

var Version = ""
var HashTag = ""
var BranchName = ""
var BuildDate = ""
var GoVersion = ""

func ShowVersion() {
	ctx := context.Background()
	ctx = outputer.SetLevel(ctx, outputer.OutputLevelNormal)
	
	outputer.ForceFL(ctx, `Version: %s
HashTag: %s
BranchName: %s
BuildDate: %s
GoVersion: %s`, Version, HashTag, BranchName, BuildDate, GoVersion)
}
