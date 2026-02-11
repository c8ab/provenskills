package integration

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// projectRoot returns the root directory of the project.
func projectRoot(t *testing.T) string {
	t.Helper()
	// Walk up from the test file location to find go.mod
	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("could not find project root (go.mod)")
		}
		dir = parent
	}
}

// buildPSK compiles the psk binary for testing.
func buildPSK(t *testing.T) string {
	t.Helper()
	root := projectRoot(t)
	bin := filepath.Join(root, "psk")
	cmd := exec.Command("go", "build", "-o", bin, "./cmd/psk/")
	cmd.Dir = root
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("failed to build psk: %v\n%s", err, out)
	}
	return bin
}

// testdataDir returns the absolute path to the testdata directory.
func testdataDir(t *testing.T) string {
	t.Helper()
	return filepath.Join(projectRoot(t), "tests", "integration", "testdata")
}

// runPSK runs the psk binary with the given args and env vars.
// Returns stdout, stderr, and the exit code.
func runPSK(t *testing.T, bin string, env []string, args ...string) (outStr, errStr string, code int) {
	t.Helper()
	cmd := exec.Command(bin, args...)
	cmd.Env = append(os.Environ(), env...)

	var stdout, stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	exitCode := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			t.Fatalf("failed to run psk: %v", err)
		}
	}

	return stdout.String(), stderr.String(), exitCode
}

func TestBuildSuccessful(t *testing.T) {
	bin := buildPSK(t)
	store := t.TempDir()

	stdout, stderr, exitCode := runPSK(t, bin,
		[]string{"PSK_STORE=" + store},
		"build", filepath.Join(testdataDir(t), "valid-skill"),
		"--maintainer", "Test <test@example.com>",
	)

	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d\nstderr: %s", exitCode, stderr)
	}

	if !strings.Contains(stdout, "Built skill: valid-skill@1.0.0") {
		t.Errorf("expected stdout to contain 'Built skill: valid-skill@1.0.0', got:\n%s", stdout)
	}

	// Verify manifest.json exists in store
	manifestPath := filepath.Join(store, "valid-skill", "1.0.0", "manifest.json")
	if _, err := os.Stat(manifestPath); err != nil {
		t.Errorf("manifest.json not found at %s: %v", manifestPath, err)
	}

	// Verify SKILL.md copied to store
	skillPath := filepath.Join(store, "valid-skill", "1.0.0", "SKILL.md")
	if _, err := os.Stat(skillPath); err != nil {
		t.Errorf("SKILL.md not found at %s: %v", skillPath, err)
	}
}

func TestBuildMissingAuthor(t *testing.T) {
	bin := buildPSK(t)
	store := t.TempDir()

	_, stderr, exitCode := runPSK(t, bin,
		[]string{"PSK_STORE=" + store},
		"build", filepath.Join(testdataDir(t), "missing-author"),
		"--maintainer", "Test <test@example.com>",
	)

	if exitCode != 2 {
		t.Fatalf("expected exit code 2, got %d", exitCode)
	}

	if !strings.Contains(stderr, "metadata.author") {
		t.Errorf("expected stderr to contain 'metadata.author', got:\n%s", stderr)
	}
}

func TestBuildMissingVersion(t *testing.T) {
	bin := buildPSK(t)
	store := t.TempDir()

	_, stderr, exitCode := runPSK(t, bin,
		[]string{"PSK_STORE=" + store},
		"build", filepath.Join(testdataDir(t), "missing-version"),
		"--maintainer", "Test <test@example.com>",
	)

	if exitCode != 2 {
		t.Fatalf("expected exit code 2, got %d", exitCode)
	}

	if !strings.Contains(stderr, "metadata.version") {
		t.Errorf("expected stderr to contain 'metadata.version', got:\n%s", stderr)
	}
}

func TestBuildMissingSKILLMD(t *testing.T) {
	bin := buildPSK(t)
	store := t.TempDir()

	_, stderr, exitCode := runPSK(t, bin,
		[]string{"PSK_STORE=" + store},
		"build", filepath.Join(testdataDir(t), "nonexistent-path"),
		"--maintainer", "Test <test@example.com>",
	)

	if exitCode != 4 {
		t.Fatalf("expected exit code 4, got %d", exitCode)
	}

	if !strings.Contains(stderr, "not found") {
		t.Errorf("expected stderr to contain 'not found', got:\n%s", stderr)
	}
}

func TestBuildMissingMaintainer(t *testing.T) {
	bin := buildPSK(t)
	store := t.TempDir()

	_, stderr, exitCode := runPSK(t, bin,
		[]string{"PSK_STORE=" + store},
		"build", filepath.Join(testdataDir(t), "valid-skill"),
	)

	if exitCode != 2 {
		t.Fatalf("expected exit code 2, got %d", exitCode)
	}

	if !strings.Contains(stderr, "--maintainer") {
		t.Errorf("expected stderr to contain '--maintainer', got:\n%s", stderr)
	}
}

func TestBuildNameDirectoryMismatch(t *testing.T) {
	bin := buildPSK(t)
	store := t.TempDir()

	// Create a temp skill dir where name field differs from directory name
	tmpDir := t.TempDir()
	skillDir := filepath.Join(tmpDir, "wrong-dir-name")
	if err := os.MkdirAll(skillDir, 0o755); err != nil {
		t.Fatal(err)
	}
	skillContent := `---
name: actual-skill-name
description: A skill with mismatched directory name.
metadata:
  version: "1.0.0"
  author: "test-author"
---

# Mismatch
`
	if err := os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte(skillContent), 0o644); err != nil {
		t.Fatal(err)
	}

	_, _, exitCode := runPSK(t, bin,
		[]string{"PSK_STORE=" + store},
		"build", skillDir,
		"--maintainer", "Test <test@example.com>",
	)

	if exitCode != 2 {
		t.Fatalf("expected exit code 2, got %d", exitCode)
	}
}

func TestBuildDuplicateConflict(t *testing.T) {
	bin := buildPSK(t)
	store := t.TempDir()

	// First build should succeed
	_, _, exitCode := runPSK(t, bin,
		[]string{"PSK_STORE=" + store},
		"build", filepath.Join(testdataDir(t), "valid-skill"),
		"--maintainer", "Test <test@example.com>",
	)
	if exitCode != 0 {
		t.Fatalf("first build expected exit code 0, got %d", exitCode)
	}

	// Second build should conflict
	_, stderr, exitCode := runPSK(t, bin,
		[]string{"PSK_STORE=" + store},
		"build", filepath.Join(testdataDir(t), "valid-skill"),
		"--maintainer", "Test <test@example.com>",
	)
	if exitCode != 3 {
		t.Fatalf("second build expected exit code 3, got %d", exitCode)
	}

	if !strings.Contains(stderr, "already exists") {
		t.Errorf("expected stderr to contain 'already exists', got:\n%s", stderr)
	}
}

func TestBuildForceOverwrite(t *testing.T) {
	bin := buildPSK(t)
	store := t.TempDir()

	// First build
	_, _, exitCode := runPSK(t, bin,
		[]string{"PSK_STORE=" + store},
		"build", filepath.Join(testdataDir(t), "valid-skill"),
		"--maintainer", "Test <test@example.com>",
	)
	if exitCode != 0 {
		t.Fatalf("first build expected exit code 0, got %d", exitCode)
	}

	// Second build with --force
	_, _, exitCode = runPSK(t, bin,
		[]string{"PSK_STORE=" + store},
		"build", filepath.Join(testdataDir(t), "valid-skill"),
		"--maintainer", "Test <test@example.com>",
		"--force",
	)
	if exitCode != 0 {
		t.Fatalf("force build expected exit code 0, got %d", exitCode)
	}
}

func TestBuildJSONOutput(t *testing.T) {
	bin := buildPSK(t)
	store := t.TempDir()

	stdout, stderr, exitCode := runPSK(t, bin,
		[]string{"PSK_STORE=" + store},
		"build", filepath.Join(testdataDir(t), "valid-skill"),
		"--maintainer", "Test <test@example.com>",
		"--json",
	)

	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d\nstderr: %s", exitCode, stderr)
	}

	var result map[string]interface{}
	if err := json.Unmarshal([]byte(stdout), &result); err != nil {
		t.Fatalf("stdout is not valid JSON: %v\nstdout: %s", err, stdout)
	}

	for _, field := range []string{"name", "version", "author", "maintainer", "path"} {
		if _, ok := result[field]; !ok {
			t.Errorf("JSON output missing field %q", field)
		}
	}
}
