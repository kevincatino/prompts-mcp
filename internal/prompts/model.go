package prompts

import "fmt"

// Prompt represents a user-defined prompt template.
type Prompt struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Content     string `json:"content"`
}

// Validate ensures required fields are present.
func (p Prompt) Validate() error {
	if p.Name == "" {
		return fmt.Errorf("name is required")
	}
	if p.Content == "" {
		return fmt.Errorf("prompt content is required for %q", p.Name)
	}
	return nil
}
