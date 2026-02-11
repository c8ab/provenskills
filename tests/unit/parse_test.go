package unit

import (
	"testing"

	"github.com/c8ab/provenskills/internal/skill"
)

func TestParseValidFrontmatter(t *testing.T) {
	input := `---
name: my-skill
description: A test skill.
metadata:
  version: "1.0.0"
  author: "test-author"
---

# My Skill

Body content here.
`
	fm, err := skill.ParseFrontmatter([]byte(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fm.Name != "my-skill" {
		t.Errorf("expected name 'my-skill', got %q", fm.Name)
	}
	if fm.Description != "A test skill." {
		t.Errorf("expected description 'A test skill.', got %q", fm.Description)
	}
	if fm.Metadata.Version != "1.0.0" {
		t.Errorf("expected version '1.0.0', got %q", fm.Metadata.Version)
	}
	if fm.Metadata.Author != "test-author" {
		t.Errorf("expected author 'test-author', got %q", fm.Metadata.Author)
	}
}

func TestParseMissingOpeningDelimiter(t *testing.T) {
	input := `name: my-skill
description: A test skill.
---

Body content.
`
	_, err := skill.ParseFrontmatter([]byte(input))
	if err == nil {
		t.Fatal("expected error for missing opening delimiter, got nil")
	}
}

func TestParseMissingClosingDelimiter(t *testing.T) {
	input := `---
name: my-skill
description: A test skill.
`
	_, err := skill.ParseFrontmatter([]byte(input))
	if err == nil {
		t.Fatal("expected error for missing closing delimiter, got nil")
	}
}

func TestParseEmptyFrontmatter(t *testing.T) {
	input := `---
---

Body content.
`
	fm, err := skill.ParseFrontmatter([]byte(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Empty frontmatter should parse but have zero-value fields
	if fm.Name != "" {
		t.Errorf("expected empty name, got %q", fm.Name)
	}
}

func TestParseAdditionalDelimitersInBody(t *testing.T) {
	input := `---
name: my-skill
description: A test skill.
metadata:
  version: "1.0.0"
  author: "test-author"
---

# My Skill

Some content with --- in it.

---

More content.
`
	fm, err := skill.ParseFrontmatter([]byte(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fm.Name != "my-skill" {
		t.Errorf("expected name 'my-skill', got %q", fm.Name)
	}
}
