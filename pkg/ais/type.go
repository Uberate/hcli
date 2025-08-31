package llms

type LLM interface {
	Thinking(input string) (resp string, err error)
}
