package ais

import (
	"context"
	"fmt"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime/model"
	"github.io/uberate/hcli/pkg/utils"
	"os"
)

type VolcConfig struct {
	ApiKey string

	ThinkModel string

	CustomPrompt map[string]string
}

func NewVolcEngineAI(config VolcConfig) VolcEngineAI {
	client := arkruntime.NewClientWithApiKey(config.ApiKey)
	res := VolcEngineAI{
		client,
		"",
		defaultSysPromptZhCn}
	// 更新用户自定义的提示词。
	for k, v := range config.CustomPrompt {
		if _, ok := res.sysRolePrompt[k]; ok {
			res.sysRolePrompt[k] = v
		}
	}

	res.thinkModelId = config.ThinkModel

	if res.thinkModelId == "" {
		res.thinkModelId = os.Getenv("THINK_MODEL_ID")
	}

	return res
}

type VolcEngineAI struct {
	c *arkruntime.Client

	thinkModelId string

	sysRolePrompt map[string]string
}

func (ve VolcEngineAI) Thinking(ctx context.Context, input string) (resp string, err error) {

	res, err := ve.c.CreateChatCompletion(ctx, model.CreateChatCompletionRequest{
		Model: ve.thinkModelId,
		Messages: []*model.ChatCompletionMessage{
			{Role: messageRoleSystem, Content: &model.ChatCompletionMessageContent{StringValue: utils.SPtr(ve.sysRolePrompt[picSummaryPromptKey])}},
			{Role: messageRoleUser, Content: &model.ChatCompletionMessageContent{StringValue: &input}},
		},
		Thinking: &model.Thinking{Type: model.ThinkingTypeDisabled},
	})

	if err != nil {
		return "", err
	}

	if len(res.Choices) == 0 {
		return "", fmt.Errorf("can't found any ans from volc")
	}

	return *res.Choices[0].Message.Content.StringValue, nil
}

const (
	messageRoleSystem = "system"
	messageRoleUser   = "user"
)
