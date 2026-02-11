package integration

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestValidateValidSkill(t *testing.T) {
	bin := buildPSK(t)

	stdout, _, exitCode := runPSK(t, bin, nil,
		"validate", filepath.Join(testdataDir(t), "valid-skill"),
	)

	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", exitCode)
	}

	if !strings.Contains(stdout, "Validation passed") {
		t.Errorf("expected stdout to contain 'Validation passed', got:\n%s", stdout)
	}

	// Check that it lists the checked fields
	for _, field := range []string{"name", "description", "version", "author", "dir match"} {
		if !strings.Contains(stdout, field) {
			t.Errorf("expected stdout to contain field %q, got:\n%s", field, stdout)
		}
	}
}

func TestValidateInvalidSkill(t *testing.T) {
	bin := buildPSK(t)

	_, stderr, exitCode := runPSK(t, bin, nil,
		"validate", filepath.Join(testdataDir(t), "bad-name"),
	)

	if exitCode != 2 {
		t.Fatalf("expected exit code 2, got %d", exitCode)
	}

	if !strings.Contains(stderr, "validation failed") || !strings.Contains(stderr, "name") {
		t.Errorf("expected stderr to list validation errors, got:\n%s", stderr)
	}
}

func TestValidateMalformedYAML(t *testing.T) {
	bin := buildPSK(t)

	_, stderr, exitCode := runPSK(t, bin, nil,
		"validate", filepath.Join(testdataDir(t), "malformed-yaml"),
	)

	if exitCode != 2 {
		t.Fatalf("expected exit code 2, got %d", exitCode)
	}

	if !strings.Contains(stderr, "error") {
		t.Errorf("expected stderr to contain parse error, got:\n%s", stderr)
	}
}
