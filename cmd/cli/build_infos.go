package main

import (
	"github.com/spf13/cobra"
	"github.io/uberate/hcli/pkg/hctx"
)

var Version = ""
var HashTag = ""
var BranchName = ""
var BuildDate = ""
var GoVersion = ""

var versionFormat = "version: %s\n" +
	"hashTag: %s\n" +
	"branchName: %s\n" +
	"buildDate: %s\n" +
	"goVersion: %s\n"

func VersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version number",
		Run: func(cmd *cobra.Command, args []string) {
			hctx.Println(cmd.Context(), versionFormat, Version, HashTag, BranchName, BuildDate, GoVersion)
		},
	}
}
