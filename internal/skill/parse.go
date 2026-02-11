// Package skill handles parsing and validation of SKILL.md files.
package skill

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

// Metadata holds the metadata block from SKILL.md frontmatter.
type Metadata struct {
	Version string `yaml:"version"`
	Author  string `yaml:"author"`
}

// SkillFrontmatter represents the parsed YAML frontmatter of a SKILL.md file.
type SkillFrontmatter struct {
	Name          string   `yaml:"name"`
	Description   string   `yaml:"description"`
	License       string   `yaml:"license"`
	Compatibility string   `yaml:"compatibility"`
	Metadata      Metadata `yaml:"metadata"`
	AllowedTools  string   `yaml:"allowed-tools"`
}

// ParseFrontmatter extracts and parses YAML frontmatter from SKILL.md content.
// The frontmatter must be delimited by --- lines at the start of the file.
func ParseFrontmatter(data []byte) (SkillFrontmatter, error) {
	content := string(data)

	// Must start with ---
	if !strings.HasPrefix(strings.TrimSpace(content), "---") {
		return SkillFrontmatter{}, fmt.Errorf("missing opening --- delimiter")
	}

	// Find the opening and closing delimiters
	trimmed := strings.TrimSpace(content)
	// Skip the first ---
	rest := trimmed[3:]

	// Find the closing ---
	idx := strings.Index(rest, "\n---")
	if idx == -1 {
		return SkillFrontmatter{}, fmt.Errorf("missing closing --- delimiter")
	}

	yamlContent := rest[:idx]

	var fm SkillFrontmatter
	if err := yaml.Unmarshal([]byte(yamlContent), &fm); err != nil {
		return SkillFrontmatter{}, fmt.Errorf("failed to parse YAML frontmatter: %w", err)
	}

	return fm, nil
}
