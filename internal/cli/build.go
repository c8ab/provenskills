package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/c8ab/provenskills/internal/exitcode"
	"github.com/c8ab/provenskills/internal/skill"
	"github.com/c8ab/provenskills/internal/store"
)

// RunBuild executes the "psk build" command.
func RunBuild(args []string) int {
	// Manual arg parsing to support intermixed flags and positional args.
	// Go's flag package stops at the first non-flag argument.
	var maintainer, path string
	var force, jsonOutput bool

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--maintainer":
			if i+1 < len(args) {
				i++
				maintainer = args[i]
			}
		case "--force":
			force = true
		case "--json":
			jsonOutput = true
		default:
			if path == "" && !strings.HasPrefix(args[i], "-") {
				path = args[i]
			}
		}
	}

	if path == "" {
		fmt.Fprintln(os.Stderr, "error: path argument is required\n\nUsage: psk build <path> --maintainer <identity>")
		return exitcode.ErrValidation
	}

	if maintainer == "" {
		fmt.Fprintln(os.Stderr, "error: --maintainer flag is required\n\nUsage: psk build <path> --maintainer <identity>")
		return exitcode.ErrValidation
	}

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
		fmt.Fprintf(os.Stderr, "error: validation failed for %s\n\n", path)
		for _, e := range errs {
			fmt.Fprintf(os.Stderr, "  - %s\n", e)
		}
		return exitcode.ErrValidation
	}

	// Normalize version
	version := skill.NormalizeVersion(fm.Metadata.Version)

	// Initialize store
	s := store.New("")
	if err := s.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to initialize store: %v\n", err)
		return exitcode.ErrIO
	}

	// Check for conflict (unless --force)
	if !force && s.Exists(fm.Name, version) {
		fmt.Fprintf(os.Stderr, "error: skill %s@%s already exists in store\n\nUse --force to overwrite.\n", fm.Name, version)
		return exitcode.ErrConflict
	}

	// Build manifest
	manifest := store.Manifest{
		ManifestVersion: 1,
		Name:            fm.Name,
		Version:         version,
		Description:     fm.Description,
		Author:          fm.Metadata.Author,
		Maintainer:      maintainer,
		BuildTimestamp:  time.Now().UTC().Format(time.RFC3339),
		Contents: store.Contents{
			SkillFile: "SKILL.md",
		},
	}

	// Add to store
	destPath, err := s.Add(fm.Name, version, path, manifest, force)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return exitcode.ErrIO
	}

	// Output
	if jsonOutput {
		result := map[string]string{
			"name":       fm.Name,
			"version":    version,
			"author":     fm.Metadata.Author,
			"maintainer": maintainer,
			"path":       destPath + "/",
		}
		data, _ := json.MarshalIndent(result, "", "  ")
		fmt.Println(string(data))
	} else {
		fmt.Printf("Built skill: %s@%s\n", fm.Name, version)
		fmt.Printf("  author:     %s\n", fm.Metadata.Author)
		fmt.Printf("  maintainer: %s\n", maintainer)
		fmt.Printf("  stored:     %s/\n", destPath)
	}

	return exitcode.Success
}
