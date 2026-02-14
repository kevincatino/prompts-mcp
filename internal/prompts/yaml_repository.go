package prompts

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// FrontmatterRepository loads prompts from Markdown files with YAML frontmatter.
type FrontmatterRepository struct {
	baseDir string
}

func NewFrontmatterRepository(baseDir string) *FrontmatterRepository {
	return &FrontmatterRepository{baseDir: baseDir}
}

func (r *FrontmatterRepository) ListPrompts(ctx context.Context) ([]Prompt, error) {
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
		if !strings.EqualFold(filepath.Ext(entry.Name()), ".md") {
			continue
		}

		name := strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name()))
		prompt, err := r.loadPrompt(entry.Name(), name)
		if err != nil {
			return nil, err
		}
		promptsList = append(promptsList, prompt)
	}

	return promptsList, nil
}

func (r *FrontmatterRepository) GetPrompt(ctx context.Context, name string) (Prompt, error) {
	select {
	case <-ctx.Done():
		return Prompt{}, ctx.Err()
	default:
	}

	filename := name + ".md"
	return r.loadPrompt(filename, name)
}

func (r *FrontmatterRepository) loadPrompt(filename string, name string) (Prompt, error) {
	content, err := os.ReadFile(filepath.Join(r.baseDir, filename))
	if err != nil {
		return Prompt{}, fmt.Errorf("read %s: %w", filename, err)
	}

	frontmatter, body, err := parseFrontmatterMarkdown(content)
	if err != nil {
		return Prompt{}, fmt.Errorf("parse %s: %w", filename, err)
	}

	var raw struct {
		Description string `yaml:"description"`
	}
	if err := yaml.Unmarshal([]byte(frontmatter), &raw); err != nil {
		return Prompt{}, fmt.Errorf("parse frontmatter in %s: %w", filename, err)
	}

	p := Prompt{
		Name:        name,
		Description: strings.TrimSpace(raw.Description),
		Content:     strings.TrimSpace(body),
	}
	if err := p.Validate(); err != nil {
		return Prompt{}, fmt.Errorf("validate %s: %w", filename, err)
	}

	return p, nil
}

func parseFrontmatterMarkdown(content []byte) (frontmatter string, body string, err error) {
	normalized := strings.ReplaceAll(string(content), "\r\n", "\n")
	lines := strings.Split(normalized, "\n")

	if len(lines) < 2 || lines[0] != "---" {
		return "", "", fmt.Errorf("missing YAML frontmatter start delimiter")
	}

	endLine := -1
	for i := 1; i < len(lines); i++ {
		if lines[i] == "---" {
			endLine = i
			break
		}
	}
	if endLine == -1 {
		return "", "", fmt.Errorf("missing YAML frontmatter end delimiter")
	}

	frontmatter = strings.Join(lines[1:endLine], "\n")
	body = strings.Join(lines[endLine+1:], "\n")
	return frontmatter, body, nil
}
