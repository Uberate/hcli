package ais

import "context"

type AIs interface {
	Thinking(ctx context.Context, input string) (resp string, err error)
	GenPic(ctx context.Context, input string) (resp []byte, err error)
}

var picSummaryPromptKey = "pic_summary_prompt_key"
var picSummaryPromptZhCn = "你是一个文本内容的描述大师，根据用户的需求可以描述需要的图片。要求：" +
	"1. 图片中不可以有任何文字内容。" +
	"2. 图片中的主要元素不可过多。" +
	"3. 响应信息的长度应该不超过 300 字。" +
	"另外，请着重描述以下内容：图片的风格，核心元素。" +
	"如果用户未指定风格、内容，你需要自行裁断选择的风格。\n" +
	"输入的信息中如果存在注释，请忽略，包括：" +
	"1. markdown 文本类型中的 '+++' 区块。"

var defaultSysPromptZhCn = map[string]string{
	picSummaryPromptKey: picSummaryPromptZhCn,
}
