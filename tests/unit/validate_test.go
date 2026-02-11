package unit

import (
	"strings"
	"testing"

	"github.com/c8ab/provenskills/internal/skill"
)

func TestValidateValidName(t *testing.T) {
	fm := skill.SkillFrontmatter{
		Name:        "my-skill",
		Description: "A valid skill.",
		Metadata:    skill.Metadata{Version: "1.0.0", Author: "test-author"},
	}
	errs := skill.Validate(fm, "my-skill")
	if len(errs) != 0 {
		t.Errorf("expected no errors, got %v", errs)
	}
}

func TestValidateUppercaseName(t *testing.T) {
	fm := skill.SkillFrontmatter{
		Name:        "My-Skill",
		Description: "A valid skill.",
		Metadata:    skill.Metadata{Version: "1.0.0", Author: "test-author"},
	}
	errs := skill.Validate(fm, "My-Skill")
	found := false
	for _, e := range errs {
		if strings.Contains(e, "name") {
			found = true
		}
	}
	if !found {
		t.Error("expected name validation error for uppercase, got none")
	}
}

func TestValidateConsecutiveHyphens(t *testing.T) {
	fm := skill.SkillFrontmatter{
		Name:        "my--skill",
		Description: "A valid skill.",
		Metadata:    skill.Metadata{Version: "1.0.0", Author: "test-author"},
	}
	errs := skill.Validate(fm, "my--skill")
	found := false
	for _, e := range errs {
		if strings.Contains(e, "name") {
			found = true
		}
	}
	if !found {
		t.Error("expected name validation error for consecutive hyphens, got none")
	}
}

func TestValidateLeadingTrailingHyphens(t *testing.T) {
	fm := skill.SkillFrontmatter{
		Name:        "-my-skill-",
		Description: "A valid skill.",
		Metadata:    skill.Metadata{Version: "1.0.0", Author: "test-author"},
	}
	errs := skill.Validate(fm, "-my-skill-")
	found := false
	for _, e := range errs {
		if strings.Contains(e, "name") {
			found = true
		}
	}
	if !found {
		t.Error("expected name validation error for leading/trailing hyphens, got none")
	}
}

func TestValidateNameTooLong(t *testing.T) {
	longName := strings.Repeat("a", 65)
	fm := skill.SkillFrontmatter{
		Name:        longName,
		Description: "A valid skill.",
		Metadata:    skill.Metadata{Version: "1.0.0", Author: "test-author"},
	}
	errs := skill.Validate(fm, longName)
	found := false
	for _, e := range errs {
		if strings.Contains(e, "name") {
			found = true
		}
	}
	if !found {
		t.Error("expected name validation error for >64 chars, got none")
	}
}

func TestValidateNameDirectoryMismatch(t *testing.T) {
	fm := skill.SkillFrontmatter{
		Name:        "my-skill",
		Description: "A valid skill.",
		Metadata:    skill.Metadata{Version: "1.0.0", Author: "test-author"},
	}
	errs := skill.Validate(fm, "other-dir")
	found := false
	for _, e := range errs {
		if strings.Contains(e, "dir") || strings.Contains(e, "match") || strings.Contains(e, "directory") {
			found = true
		}
	}
	if !found {
		t.Error("expected name-directory mismatch error, got none")
	}
}

func TestValidateEmptyDescription(t *testing.T) {
	fm := skill.SkillFrontmatter{
		Name:        "my-skill",
		Description: "",
		Metadata:    skill.Metadata{Version: "1.0.0", Author: "test-author"},
	}
	errs := skill.Validate(fm, "my-skill")
	found := false
	for _, e := range errs {
		if strings.Contains(e, "description") {
			found = true
		}
	}
	if !found {
		t.Error("expected description validation error, got none")
	}
}

func TestValidateDescriptionTooLong(t *testing.T) {
	fm := skill.SkillFrontmatter{
		Name:        "my-skill",
		Description: strings.Repeat("a", 1025),
		Metadata:    skill.Metadata{Version: "1.0.0", Author: "test-author"},
	}
	errs := skill.Validate(fm, "my-skill")
	found := false
	for _, e := range errs {
		if strings.Contains(e, "description") {
			found = true
		}
	}
	if !found {
		t.Error("expected description too long error, got none")
	}
}

func TestValidateValidSemver(t *testing.T) {
	fm := skill.SkillFrontmatter{
		Name:        "my-skill",
		Description: "A valid skill.",
		Metadata:    skill.Metadata{Version: "2.1.3", Author: "test-author"},
	}
	errs := skill.Validate(fm, "my-skill")
	if len(errs) != 0 {
		t.Errorf("expected no errors, got %v", errs)
	}
}

func TestValidateMajorMinorNormalization(t *testing.T) {
	fm := skill.SkillFrontmatter{
		Name:        "my-skill",
		Description: "A valid skill.",
		Metadata:    skill.Metadata{Version: "1.0", Author: "test-author"},
	}
	errs := skill.Validate(fm, "my-skill")
	if len(errs) != 0 {
		t.Errorf("expected no errors for major.minor format, got %v", errs)
	}
}

func TestValidateNormalizeVersion(t *testing.T) {
	v := skill.NormalizeVersion("1.0")
	if v != "1.0.0" {
		t.Errorf("expected '1.0.0', got %q", v)
	}
	v2 := skill.NormalizeVersion("2.1.3")
	if v2 != "2.1.3" {
		t.Errorf("expected '2.1.3', got %q", v2)
	}
}

func TestValidateInvalidVersion(t *testing.T) {
	fm := skill.SkillFrontmatter{
		Name:        "my-skill",
		Description: "A valid skill.",
		Metadata:    skill.Metadata{Version: "abc", Author: "test-author"},
	}
	errs := skill.Validate(fm, "my-skill")
	found := false
	for _, e := range errs {
		if strings.Contains(e, "version") {
			found = true
		}
	}
	if !found {
		t.Error("expected version validation error, got none")
	}
}

func TestValidateEmptyAuthor(t *testing.T) {
	fm := skill.SkillFrontmatter{
		Name:        "my-skill",
		Description: "A valid skill.",
		Metadata:    skill.Metadata{Version: "1.0.0", Author: ""},
	}
	errs := skill.Validate(fm, "my-skill")
	found := false
	for _, e := range errs {
		if strings.Contains(e, "author") {
			found = true
		}
	}
	if !found {
		t.Error("expected author validation error, got none")
	}
}
