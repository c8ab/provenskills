// Package store handles the local Proven Skill Artifact store.
package store

import (
	"encoding/json"
	"fmt"
	"os"
)

// Contents describes the files in a stored skill artifact.
type Contents struct {
	SkillFile  string   `json:"skillFile"`
	Scripts    []string `json:"scripts,omitempty"`
	References []string `json:"references,omitempty"`
	Assets     []string `json:"assets,omitempty"`
}

// Manifest represents the metadata for a stored skill artifact.
type Manifest struct {
	ManifestVersion int      `json:"manifestVersion"`
	Name            string   `json:"name"`
	Version         string   `json:"version"`
	Description     string   `json:"description"`
	Author          string   `json:"author"`
	Maintainer      string   `json:"maintainer"`
	BuildTimestamp  string   `json:"buildTimestamp"`
	Contents        Contents `json:"contents,omitempty"`
	SourceHash      string   `json:"sourceHash,omitempty"`
}

// WriteManifest writes a Manifest to the given file path as JSON.
func WriteManifest(path string, m Manifest) error {
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal manifest: %w", err)
	}
	data = append(data, '\n')
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("failed to write manifest: %w", err)
	}
	return nil
}

// ReadManifest reads and parses a manifest.json file.
func ReadManifest(path string) (Manifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Manifest{}, fmt.Errorf("failed to read manifest: %w", err)
	}
	var m Manifest
	if err := json.Unmarshal(data, &m); err != nil {
		return Manifest{}, fmt.Errorf("failed to parse manifest: %w", err)
	}
	return m, nil
}
