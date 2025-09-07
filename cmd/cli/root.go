package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.io/uberate/hcli/cmd/cli/cmds"
	"github.io/uberate/hcli/pkg/ais"
	"github.io/uberate/hcli/pkg/cctx"
	"github.io/uberate/hcli/pkg/config"
	"github.io/uberate/hcli/pkg/outputer"
	"os"
	"path"
	"path/filepath"
)

var showVersion bool
var showDetail bool
var silence bool
var configPath string

func RootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Run: func(cmd *cobra.Command, args []string) {
			return
		},
	}

	cmd.PersistentFlags().BoolVarP(&showVersion, "version", "v", false, "show version, exit directly")
	cmd.PersistentFlags().BoolVarP(&showDetail, "detail", "d", false, "show detail, conflict with silence mode")
	cmd.PersistentFlags().BoolVarP(&silence, "silence", "s", false, "silence mode, conflict with detail mode")
	cmd.PersistentFlags().StringVarP(&configPath, "config", "c", "", "config path, "+
		"if empty, search config by these path, './.hcli_config.yaml', '~/.hcli_config.yaml', "+
		"'{parents_dir}/.hcli_config.yaml")

	cmd.AddCommand(
		cmds.GenCmd(),
		cmds.ConfigCmd(),
	)

	cmd.PersistentPreRunE = preRun
	return cmd
}

func preRun(cmd *cobra.Command, args []string) error {

	cPath := cmd.CommandPath()
	if cPath == " config demo-config" {
		return nil
	}

	ctx := cmd.Context()

	if showDetail && silence {
		return errors.New("can not use both --silence and --detail")
	}

	if showDetail {
		ctx = outputer.SetLevel(ctx, outputer.OutputLevelDetail)
		outputer.DetailFL(ctx, "enable detail mode")
	} else if silence {
		ctx = outputer.SetLevel(ctx, outputer.OutputLevelSilence)
	} else {
		ctx = outputer.SetLevel(ctx, outputer.OutputLevelNormal)
	}

	if showVersion {
		ShowVersion()
		return nil
	}

	c, err := loadConfig(ctx)
	if err != nil {
		return err
	}

	ctx = config.SetConfig(ctx, c)

	// Create and set AI client if needed
	aiClient := createAIClient(c)
	ctx = cctx.SetAIClient(ctx, aiClient)

	cmd.SetContext(ctx)

	return nil
}

func loadConfig(ctx context.Context) (config.CliConfig, error) {
	if configPath != "" {
		outputer.DetailFL(ctx, "read config from: %s", configPath)
		return config.ReadConfig(configPath)
	}

	outputer.DetailFL(ctx, "config path not set, search './.hcli_config.yaml'")
	// 目录不存在按照以下顺序读取。
	res, err := config.ReadConfig("./.hcli_config.yaml")
	if !os.IsNotExist(err) {
		return res, err
	}

	outputer.DetailFL(ctx, "config not found: './.hcli_config.yaml'")
	outputer.DetailFL(ctx, "search: '~/.hcli_config.yaml'")
	res, err = config.ReadConfig("~/.hcli_config.yaml")
	if !os.IsNotExist(err) {
		return res, err
	}

	outputer.DetailFL(ctx, "~/.hcli_config.yaml not found, search parent dirs")
	execPath, err := os.Executable()
	if err != nil {
		outputer.DetailFL(ctx, "search current path fail: %v", err)
		return res, err
	}

	execPath = filepath.Dir(execPath) // get dir

	for ; !(execPath == "" || execPath == "/"); execPath = filepath.Dir(execPath) {
		readPath := path.Join(execPath, ".hcli_config.yaml")
		outputer.DetailFL(ctx, "try read path: %s", readPath)
		res, err = config.ReadConfig(readPath)
		if !os.IsNotExist(err) {
			outputer.DetailFL(ctx, "try read path: %s fail: %v", readPath, err)
			return res, err
		}
	}

	return res, fmt.Errorf("nout found config, please set config to .hcli_config.yaml by hcli config demo")
}

func createAIClient(cfg config.CliConfig) ais.AIs {
	aiConfig := cfg.AI
	
	switch aiConfig.Provider {
	case "volc":
		return createVolcEngineAI(aiConfig)
	// Add more AI providers here in the future
	// case "openai":
	//     return createOpenAI(aiConfig)
	// case "anthropic":
	//     return createAnthropic(aiConfig)
	default:
		// Default to VolcEngine if provider is not specified or unknown
		return createVolcEngineAI(aiConfig)
	}
}

func createVolcEngineAI(aiConfig config.AIConfig) ais.AIs {
	// Use config values first, fall back to environment variables
	apiKey := aiConfig.APIKey
	if apiKey == "" {
		apiKey = os.Getenv("VOLC_API_KEY")
	}
	
	thinkModel := aiConfig.ThinkModel
	if thinkModel == "" {
		thinkModel = os.Getenv("THINK_MODEL_ID")
	}
	
	picModel := aiConfig.PicModel
	if picModel == "" {
		picModel = os.Getenv("PIC_MODEL_ID")
	}
	
	// Use custom prompts from config or default empty map
	customPrompt := aiConfig.CustomPrompt
	if customPrompt == nil {
		customPrompt = make(map[string]string)
	}
	
	volcConfig := ais.VolcConfig{
		ApiKey:       apiKey,
		ThinkModel:   thinkModel,
		PicModel:     picModel,
		CustomPrompt: customPrompt,
	}
	
	return ais.NewVolcEngineAI(volcConfig)
}
