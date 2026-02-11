package cli

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/c8ab/provenskills/internal/exitcode"
	"github.com/c8ab/provenskills/internal/store"
)

// RunList executes the "psk list" command.
func RunList(args []string) int {
	fs := flag.NewFlagSet("list", flag.ContinueOnError)
	jsonOutput := fs.Bool("json", false, "Output as JSON array")
	fs.SetOutput(os.Stderr)

	if err := fs.Parse(args); err != nil {
		return exitcode.ErrValidation
	}

	s := store.New("")
	manifests, err := s.List()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return exitcode.ErrIO
	}

	if *jsonOutput {
		if len(manifests) == 0 {
			fmt.Println("[]")
			return exitcode.Success
		}
		type listEntry struct {
			Name        string `json:"name"`
			Version     string `json:"version"`
			Description string `json:"description"`
			Author      string `json:"author"`
			Maintainer  string `json:"maintainer"`
		}
		var entries []listEntry
		for _, m := range manifests {
			entries = append(entries, listEntry{
				Name:        m.Name,
				Version:     m.Version,
				Description: m.Description,
				Author:      m.Author,
				Maintainer:  m.Maintainer,
			})
		}
		data, _ := json.MarshalIndent(entries, "", "  ")
		fmt.Println(string(data))
		return exitcode.Success
	}

	if len(manifests) == 0 {
		fmt.Println("No skills found in store.")
		return exitcode.Success
	}

	// Calculate column widths
	nameW, versionW, authorW := 4, 7, 6 // header lengths
	for _, m := range manifests {
		if len(m.Name) > nameW {
			nameW = len(m.Name)
		}
		if len(m.Version) > versionW {
			versionW = len(m.Version)
		}
		if len(m.Author) > authorW {
			authorW = len(m.Author)
		}
	}

	fmtStr := fmt.Sprintf("%%-%ds  %%-%ds  %%-%ds  %%s\n", nameW, versionW, authorW)
	fmt.Printf(fmtStr, "NAME", "VERSION", "AUTHOR", "MAINTAINER")
	for _, m := range manifests {
		fmt.Printf(fmtStr, m.Name, m.Version, m.Author, m.Maintainer)
	}

	return exitcode.Success
}
