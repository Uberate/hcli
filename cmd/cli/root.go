package main

import (
	"github.com/spf13/cobra"
	"github.io/uberate/hcli/cmd/cli/cmds"
	"github.io/uberate/hcli/pkg/config"
	"github.io/uberate/hcli/pkg/hctx"
	"github.io/uberate/hcli/pkg/output"
)

var configPath string
var logLevel string

func RootCmd() *cobra.Command {

	cmd := &cobra.Command{
		Run: func(cmd *cobra.Command, args []string) {
			hctx.Println(cmd.Context(), "hcli")
			return
		},
	}

	cmd.AddCommand(
		VersionCmd(),
		cmds.ConfigCmd(),
		cmds.GenCmd(),
	)
	cmd.PersistentFlags().StringVarP(&logLevel, "log-level", "l", "info", "log level, support: debug, info, warn, error, fatal")
	cmd.PersistentFlags().StringVarP(&configPath, "config", "c", ".hcli_config.yaml", "config file path,"+
		" if empty, use './.hcli_config.yaml''")

	cmd.PreRunE = preRun

	return cmd
}

func preRun(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()

	l, err := output.ParseLevel(logLevel)
	if err != nil {
		return err
	}

	ctx = hctx.SetOutputter(ctx, output.NewOutputter(l, cmd.OutOrStdout()))
	ctx = hctx.SetConfigPath(ctx, configPath)
	cmd.SetContext(ctx)
	return nil
}

func readConfig() (config.CliConfig, error) {
	if configPath == "" {
		configPath = "./.hcli_config.yaml"
	}

	return config.ReadConfig(configPath)
}
