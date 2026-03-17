package tui

type Category struct {
	Name  string
	Query string
}

var Categories = []Category{
	{Name: "Claude / Anthropic", Query: "claude OR anthropic OR mcp-server topic:claude"},
	{Name: "AI / LLM", Query: "llm OR ai OR gpt topic:machine-learning"},
	{Name: "Go", Query: "language:go stars:>1000"},
	{Name: "Web", Query: "topic:react OR topic:vue OR topic:nextjs"},
}
