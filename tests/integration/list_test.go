package integration

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestListWithSkills(t *testing.T) {
	bin := buildPSK(t)
	store := t.TempDir()

	// Build two skills
	_, _, exitCode := runPSK(t, bin,
		[]string{"PSK_STORE=" + store},
		"build", filepath.Join(testdataDir(t), "valid-skill"),
		"--maintainer", "Test <test@example.com>",
	)
	if exitCode != 0 {
		t.Fatalf("build valid-skill failed with exit code %d", exitCode)
	}

	// Build a second skill (using missing-version fixture with a fixed version for list test)
	// We need a second valid skill - create one on the fly
	skillDir := createTempSkill(t, "another-skill", "2.0.0", "other-author")
	_, _, exitCode = runPSK(t, bin,
		[]string{"PSK_STORE=" + store},
		"build", skillDir,
		"--maintainer", "Other <other@example.com>",
	)
	if exitCode != 0 {
		t.Fatalf("build another-skill failed with exit code %d", exitCode)
	}

	// List skills
	stdout, _, exitCode := runPSK(t, bin,
		[]string{"PSK_STORE=" + store},
		"list",
	)
	if exitCode != 0 {
		t.Fatalf("list expected exit code 0, got %d", exitCode)
	}

	if !strings.Contains(stdout, "valid-skill") {
		t.Errorf("expected stdout to contain 'valid-skill', got:\n%s", stdout)
	}
	if !strings.Contains(stdout, "another-skill") {
		t.Errorf("expected stdout to contain 'another-skill', got:\n%s", stdout)
	}
	if !strings.Contains(stdout, "1.0.0") {
		t.Errorf("expected stdout to contain '1.0.0', got:\n%s", stdout)
	}
	if !strings.Contains(stdout, "2.0.0") {
		t.Errorf("expected stdout to contain '2.0.0', got:\n%s", stdout)
	}
	if !strings.Contains(stdout, "test-author") {
		t.Errorf("expected stdout to contain 'test-author', got:\n%s", stdout)
	}
	if !strings.Contains(stdout, "Test <test@example.com>") {
		t.Errorf("expected stdout to contain maintainer, got:\n%s", stdout)
	}
}

func TestListEmptyStore(t *testing.T) {
	bin := buildPSK(t)
	store := t.TempDir()

	stdout, _, exitCode := runPSK(t, bin,
		[]string{"PSK_STORE=" + store},
		"list",
	)
	if exitCode != 0 {
		t.Fatalf("list expected exit code 0, got %d", exitCode)
	}

	if !strings.Contains(stdout, "No skills found") {
		t.Errorf("expected stdout to contain 'No skills found', got:\n%s", stdout)
	}
}

func TestListJSON(t *testing.T) {
	bin := buildPSK(t)
	store := t.TempDir()

	// Build a skill
	_, _, exitCode := runPSK(t, bin,
		[]string{"PSK_STORE=" + store},
		"build", filepath.Join(testdataDir(t), "valid-skill"),
		"--maintainer", "Test <test@example.com>",
	)
	if exitCode != 0 {
		t.Fatalf("build failed with exit code %d", exitCode)
	}

	stdout, _, exitCode := runPSK(t, bin,
		[]string{"PSK_STORE=" + store},
		"list", "--json",
	)
	if exitCode != 0 {
		t.Fatalf("list --json expected exit code 0, got %d", exitCode)
	}

	var result []map[string]interface{}
	if err := json.Unmarshal([]byte(stdout), &result); err != nil {
		t.Fatalf("stdout is not valid JSON: %v\nstdout: %s", err, stdout)
	}

	if len(result) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(result))
	}

	for _, field := range []string{"name", "version", "description", "author", "maintainer"} {
		if _, ok := result[0][field]; !ok {
			t.Errorf("JSON output missing field %q", field)
		}
	}
}

// createTempSkill creates a temporary skill directory for testing.
func createTempSkill(t *testing.T, name, version, author string) string {
	t.Helper()
	tmpDir := t.TempDir()
	skillDir := filepath.Join(tmpDir, name)
	if err := os.MkdirAll(skillDir, 0o755); err != nil {
		t.Fatal(err)
	}
	content := fmt.Sprintf(`---
name: %s
description: A test skill.
metadata:
  version: "%s"
  author: "%s"
---

# %s
`, name, version, author, name)
	if err := os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return skillDir
}
