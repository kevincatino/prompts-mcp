package prompts

import "context"

// Repository provides access to prompt definitions.
type Repository interface {
	ListPrompts(ctx context.Context) ([]Prompt, error)
	GetPrompt(ctx context.Context, name string) (Prompt, error)
}
