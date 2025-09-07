# AI Agent Configuration

This is the description of `hcli` project.
`hcli` written by `go` and use `cobra` as command framework.

Hcli was a tool to quick complete the hugo posts.

- [x] Support `hcli gen posts -n template` to generate a new posts.
- [ ] Support read the posts to create a feature.png for some docs. 
      The commands like `hcli gen feature -n template`.
- [ ] Support models: volc(Vol):
  - Chat API: https://www.volcengine.com/docs/82379/1494384
  - Text to PNG API: https://www.volcengine.com/docs/82379/1541523
## Config demo
```yaml
# Define the template of posts.
Templates:
  - Name: sa # hcli gen posts -nsa xxx
    Categories: ["阅读"]
    Tags: ["短篇阅读"]
    Template: |
      +++
      date = '{{.createAt}}'
      title = '{{.title}}'
      categories = {{.categories}}
      tags = {{.tags}}
      +++

      > 原文 [Title](Url) -- Author
    NeedDir: true
    Dir: "content/posts/article/short"
    PicUserPrompt: "" # use default prompt, only need dir was true can generate pic.
```

## Code format

You must follow this format to write code:

1. All commands was define at `cmd/cli.cmds`.
For any command non-buss action was define here too.
Like: 
    - Action param check.
    - Init some file.
    - Read the config.

2. All buss code was define in pkg.

## Dev environment

Use 