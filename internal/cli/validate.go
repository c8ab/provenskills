package cli

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/c8ab/provenskills/internal/exitcode"
	"github.com/c8ab/provenskills/internal/skill"
)

// RunValidate executes the "psk validate" command.
func RunValidate(args []string) int {
	fs := flag.NewFlagSet("validate", flag.ContinueOnError)
	jsonOutput := fs.Bool("json", false, "Output result as JSON")
	fs.SetOutput(os.Stderr)

	if err := fs.Parse(args); err != nil {
		return exitcode.ErrValidation
	}

	if fs.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "error: path argument is required\n\nUsage: psk validate <path>")
		return exitcode.ErrValidation
	}

	path := fs.Arg(0)

	// Reject path traversal
	if containsPathTraversal(path) {
		fmt.Fprintln(os.Stderr, "error: path contains '..' segments (path traversal not allowed)")
		return exitcode.ErrValidation
	}

	// Check path exists
	skillMDPath := filepath.Join(path, "SKILL.md")
	if _, err := os.Stat(skillMDPath); err != nil {
		if os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "error: directory not found: %s\n", path)
		} else {
			fmt.Fprintf(os.Stderr, "error: cannot access %s: %v\n", path, err)
		}
		return exitcode.ErrIO
	}

	// Read and parse SKILL.md
	data, err := os.ReadFile(skillMDPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to read SKILL.md: %v\n", err)
		return exitcode.ErrIO
	}

	fm, err := skill.ParseFrontmatter(data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: validation failed for %s\n\n  - %v\n", path, err)
		return exitcode.ErrValidation
	}

	// Validate
	dirName := filepath.Base(path)
	errs := skill.Validate(fm, dirName)

	if len(errs) > 0 {
		if *jsonOutput {
			result := map[string]interface{}{
				"valid":  false,
				"path":   path,
				"errors": errs,
			}
			out, _ := json.MarshalIndent(result, "", "  ")
			fmt.Fprintln(os.Stderr, string(out))
		} else {
			fmt.Fprintf(os.Stderr, "error: validation failed for %s\n\n", path)
			for _, e := range errs {
				fmt.Fprintf(os.Stderr, "  - %s\n", e)
			}
		}
		return exitcode.ErrValidation
	}

	// Success output
	version := skill.NormalizeVersion(fm.Metadata.Version)

	if *jsonOutput {
		result := map[string]interface{}{
			"valid":       true,
			"path":        path,
			"name":        fm.Name,
			"description": fmt.Sprintf("present (%d chars)", len(fm.Description)),
			"version":     version,
			"author":      fm.Metadata.Author,
			"dirMatch":    true,
		}
		out, _ := json.MarshalIndent(result, "", "  ")
		fmt.Println(string(out))
	} else {
		fmt.Printf("Validation passed: %s\n\n", path)
		fmt.Printf("  name:        %s (valid)\n", fm.Name)
		fmt.Printf("  description: present (%d chars)\n", len(fm.Description))
		fmt.Printf("  version:     %s (valid semver)\n", version)
		fmt.Printf("  author:      %s (present)\n", fm.Metadata.Author)
		fmt.Printf("  dir match:   %s == %s (ok)\n", fm.Name, dirName)
	}

	return exitcode.Success
}
