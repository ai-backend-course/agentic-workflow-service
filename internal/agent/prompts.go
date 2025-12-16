package agent

import (
	"bytes"
	"embed"
	"fmt"
	"strings"
	"text/template"
)

//go:embed *.txt
var promptFS embed.FS

type PromptLoader interface {
	Render(name string, data map[string]any) (string, error)
}

type EmbeddedPromptLoader struct {
	templates map[string]*template.Template
}

func NewEmbeddedPromptLoader() (*EmbeddedPromptLoader, error) {
	loader := &EmbeddedPromptLoader{
		templates: make(map[string]*template.Template),
	}

	entries, err := promptFS.ReadDir(".")
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		content, err := promptFS.ReadFile(name)
		if err != nil {
			return nil, err
		}

		key := strings.TrimSuffix(name, ".txt")

		tmpl, err := template.New(key).Parse(string(content))
		if err != nil {
			return nil, err
		}

		loader.templates[key] = tmpl
	}

	if len(loader.templates) == 0 {
		return nil, fmt.Errorf("no embedded prompts found")
	}

	return loader, nil
}

func (pl *EmbeddedPromptLoader) Render(name string, data map[string]any) (string, error) {
	tmpl, ok := pl.templates[name]
	if !ok {
		return "", fmt.Errorf("prompt not found: %s", name)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
