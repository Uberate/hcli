package llms

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime/model"
)

type VolcEngineConfig struct {
	ApiKey    string `yaml:"ApiKey" describe:"The API Key to access volc engine"`
	TextModel string `yaml:"TextModel" describe:"The TextModel id"`
	PicModel  string `yaml:"PicModel" describe:"The PicModel id"`
}

func NewVolcEngineLLM(vec VolcEngineConfig) LLMTools {
	client := arkruntime.NewClientWithApiKey(vec.ApiKey)

	res := &VolcEngineLLM{
		client,
		vec.TextModel,
		vec.PicModel,
	}

	return res
}

type VolcEngineLLM struct {
	c *arkruntime.Client

	textModel string
	picModel  string
}

func (vel *VolcEngineLLM) Text(ctx context.Context, sysPrompt, input string) (resp string, err error) {
	res, err := vel.c.CreateChatCompletion(ctx, model.CreateChatCompletionRequest{
		Model: vel.textModel,
		Messages: []*model.ChatCompletionMessage{
			{Role: messageRoleSystem, Content: &model.ChatCompletionMessageContent{StringValue: &sysPrompt}},
			{Role: messageRoleUser, Content: &model.ChatCompletionMessageContent{StringValue: &input}},
		},
	})

	if err != nil {
		return "", err
	}

	if len(res.Choices) == 0 {
		return "", fmt.Errorf("no response choices received from VolcEngine LLM - check API key and model configuration")
	}

	return *res.Choices[0].Message.Content.StringValue, nil
}

func (vel *VolcEngineLLM) Pic(ctx context.Context, input string) (resp []byte, err error) {
	if vel.picModel == "" {
		return nil, fmt.Errorf("no pic model specified for VolcEngine LLM")
	}

	form := model.GenerateImagesResponseFormatBase64
	result, err := vel.c.GenerateImages(ctx, model.GenerateImagesRequest{
		Model:          vel.picModel,
		Prompt:         input,
		ResponseFormat: &form,
	})

	if err != nil {
		return nil, err
	}

	if len(result.Data) == 0 {
		return nil, fmt.Errorf("no image data received from VolcEngine LLM - check API key and model configuration")
	}

	if result.Data[0].B64Json == nil {
		return nil, fmt.Errorf("iamge data format invalid - expected base64 encoded image but received nil")
	}

	// Decode base64 image data
	imageData, err := base64.StdEncoding.DecodeString(*result.Data[0].B64Json)
	if err != nil {
		return nil, fmt.Errorf("failed to decode nase64 image data: %w", err)
	}

	return imageData, nil
}

const (
	messageRoleSystem = "system"
	messageRoleUser   = "user"
)
