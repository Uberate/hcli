package config

type TemplateConfig struct {
	Name string `yaml:"Name" comment:"Template name, use hcli gen posts --template-name(|-n)=Name to generate a new posts
by template value to init."`

	Categories []string `yaml:"Categories" comment:"Categories of posts"`

	Tags []string `yaml:"Tags" comment:"Tags of posts"`

	Template string `yaml:"Template" comment:"The go template of posts. Default template was: \n
+++
todo
+++
Variables:
- 
TODO"`

	Dir     string `yaml:"Dir" comment:"Generate post path"`
	NeedDir bool   `yaml:"NeedDir" comment:"if need dir, hcli will create post in a new dir named args and set file to index.md" default:"false"`

	PicSummaryPrompt string `yaml:"PicSummaryPrompt" comment:"Pic summary prompt from target file."`
	PicCreatePrompt  string `yaml:"PicCreatePrompt" comment:"Pic create prompt from target file."`
}
