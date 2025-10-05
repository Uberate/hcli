package poster

import (
	"context"
	"errors"
	"github.io/uberate/hcli/pkg/hctx"
	"github.io/uberate/hcli/pkg/llms"
	"github.io/uberate/hcli/pkg/template"
)

type GeneratePosterArgs struct {
	TP       *template.Template
	LLMTools llms.LLMTools
	FileName string
}

type GeneratePosterResult struct {
	Summary string
	Pic     []byte
}

func GeneratePoster(ctx context.Context, args GeneratePosterArgs) (*GeneratePosterResult, error) {
	if args.TP == nil {
		return nil, errors.New("template is required for generate poster")
	}

	if args.LLMTools == nil {
		return nil, errors.New("LLMTools is required for generate poster")
	}

	fileContent, err := args.TP.ReadIfExists(args.FileName)
	if err != nil {
		return nil, err
	}

	summaryPrompt := args.TP.PicSummaryPrompt
	if summaryPrompt == "" {
		summaryPrompt = defaultPicSummaryPrompt
	}

	res, err := args.LLMTools.Text(ctx, summaryPrompt, string(fileContent))
	if err != nil {
		return nil, err
	}
	hctx.Println(ctx, "generate the summary done")

	picData, err := args.LLMTools.Pic(ctx, res)
	if err != nil {
		return nil, err
	}

	hctx.Println(ctx, "generate the pic done")

	return &GeneratePosterResult{
		Summary: res,
		Pic:     picData,
	}, nil

}

var defaultPicSummaryPrompt = "" +
	"你是一个文本内容的描述大师，根据用户的需求可以描述需要的图片。要求：" +
	"1. 图片中不可以有任何文字内容。" +
	"2. 图片中的主要元素不可过多。" +
	"3. 响应信息的长度应该不超过 300 字。" +
	"另外，请着重描述以下内容：图片的风格，核心元素。" +
	"如果用户未指定风格、内容，你需要自行裁断选择的风格。\n" +
	"输入的信息中如果存在注释，请忽略：" +
	"1. markdown 文本类型中的 '+++' 区块。"
