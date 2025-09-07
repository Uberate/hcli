package ais

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime/model"
	"github.io/uberate/hcli/pkg/utils"
	"os"
)

type VolcConfig struct {
	ApiKey string

	ThinkModel string
	PicModel   string

	CustomPrompt map[string]string
}

func NewVolcEngineAI(config VolcConfig) VolcEngineAI {
	client := arkruntime.NewClientWithApiKey(config.ApiKey)
	res := VolcEngineAI{
		client,
		config.ThinkModel,
		config.PicModel,
		defaultSysPromptZhCn}
	// 更新用户自定义的提示词。
	for k, v := range config.CustomPrompt {
		if _, ok := res.sysRolePrompt[k]; ok {
			res.sysRolePrompt[k] = v
		}
	}

	res.thinkModelId = config.ThinkModel
	res.picModelId = config.PicModel

	if res.thinkModelId == "" {
		res.thinkModelId = os.Getenv("THINK_MODEL_ID")
	}
	if res.picModelId == "" {
		res.picModelId = os.Getenv("PIC_MODEL_ID")
	}

	return res
}

type VolcEngineAI struct {
	c *arkruntime.Client

	thinkModelId string
	picModelId   string

	sysRolePrompt map[string]string
}

func (ve VolcEngineAI) CreatePICSummary(ctx context.Context, input string) (resp string, err error) {

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

func (ve VolcEngineAI) GenPic(ctx context.Context, input string) (resp []byte, err error) {
	if ve.picModelId == "" {
		return nil, fmt.Errorf("pic model not configured")
	}

	form := model.GenerateImagesResponseFormatBase64
	result, err := ve.c.GenerateImages(ctx, model.GenerateImagesRequest{
		Model:          ve.picModelId,
		Prompt:         input,
		ResponseFormat: &form,
	})

	if err != nil {
		return nil, err
	}

	if len(result.Data) == 0 {
		return nil, fmt.Errorf("can't found any ans from volc")
	}

	if result.Data[0].B64Json == nil {
		return nil, fmt.Errorf("can't found any ans from volc")
	}

	// Decode base64 image data
	imageData, err := base64.StdEncoding.DecodeString(*result.Data[0].B64Json)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64 image data: %w", err)
	}

	return imageData, nil
}

const (
	messageRoleSystem = "system"
	messageRoleUser   = "user"
)
