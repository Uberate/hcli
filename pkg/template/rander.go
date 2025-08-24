package template

import (
	"bytes"
	"context"
	"fmt"
	"github.io/uberate/hcli/pkg/outputer"
	"strings"
	"text/template"
	"time"
)

type RenderConfig struct {
	Tags       []string
	Categories []string

	Time string

	CustomArgs map[string]string // can cover exists args
	Temp       string

	Title string
}

func Render(ctx context.Context, cc RenderConfig) (string, error) {

	vars := map[string]string{
		"title":      cc.Title,
		"createAt":   cc.Time,
		"tags":       formatStringArray(cc.Tags),
		"categories": formatStringArray(cc.Categories),
	}

	for k, v := range cc.CustomArgs {
		vars[k] = v
	}

	for k, v := range vars {
		outputer.DetailFL(ctx, "Var: %s=%s", k, v)
	}

	body := defaultTemplate
	if cc.Temp != "" {
		body = cc.Temp
	}

	temp, err := template.New(cc.Title).Parse(body)
	if err != nil {
		return "", err
	}

	res := bytes.NewBuffer(nil)

	err = temp.Execute(res, vars)
	return res.String(), err
}

func formatStringArray(input []string) string {
	res := strings.Builder{}
	res.WriteString("[")
	for index, item := range input {
		if index != 0 {
			res.WriteString(",")
		}

		res.WriteString(fmt.Sprintf("\"%s\"", item))
	}

	res.WriteString("]")
	return res.String()
}

func tryParseTimeLayout(ctx context.Context, input string) string {
	switch strings.ToLower(input) {
	case "rfc3339":
		return time.RFC3339
	case "rfc3339nano":
		return time.RFC3339Nano
	case "ansic":
		return time.ANSIC
	case "layout":
		return time.Layout
	}

	return input
}

func getTime(ctx context.Context, t time.Time, format string) string {

	f := tryParseTimeLayout(ctx, format)
	outputer.DetailFL(ctx, "try to render time by format: %s", f)
	return t.Format(f)
}
