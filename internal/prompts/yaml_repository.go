package prompts

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// YAMLRepository loads prompts from YAML files in a directory.
type YAMLRepository struct {
	baseDir string
}

func NewYAMLRepository(baseDir string) *YAMLRepository {
	return &YAMLRepository{baseDir: baseDir}
}

func (r *YAMLRepository) ListPrompts(ctx context.Context) ([]Prompt, error) {
	entries, err := os.ReadDir(r.baseDir)
	if err != nil {
		return nil, fmt.Errorf("read prompts dir: %w", err)
	}

	var promptsList []Prompt
	for _, entry := range entries {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		if entry.IsDir() {
			continue
		}
		if filepath.Ext(entry.Name()) != ".yaml" {
			continue
		}

		name := strings.TrimSuffix(entry.Name(), ".yaml")
		prompt, err := r.loadPrompt(entry.Name(), name)
		if err != nil {
			return nil, err
		}
		promptsList = append(promptsList, prompt)
	}

	return promptsList, nil
}

func (r *YAMLRepository) GetPrompt(ctx context.Context, name string) (Prompt, error) {
	select {
	case <-ctx.Done():
		return Prompt{}, ctx.Err()
	default:
	}

	filename := name + ".yaml"
	return r.loadPrompt(filename, name)
}

func (r *YAMLRepository) loadPrompt(filename string, name string) (Prompt, error) {
	content, err := os.ReadFile(filepath.Join(r.baseDir, filename))
	if err != nil {
		return Prompt{}, fmt.Errorf("read %s: %w", filename, err)
	}

	var raw struct {
		Description string `yaml:"description"`
		Prompt      string `yaml:"prompt"`
	}
	if err := yaml.Unmarshal(content, &raw); err != nil {
		return Prompt{}, fmt.Errorf("parse %s: %w", filename, err)
	}

	p := Prompt{
		Name:        name,
		Description: strings.TrimSpace(raw.Description),
		Content:     strings.TrimSpace(raw.Prompt),
	}
	if err := p.Validate(); err != nil {
		return Prompt{}, fmt.Errorf("validate %s: %w", filename, err)
	}

	return p, nil
}
