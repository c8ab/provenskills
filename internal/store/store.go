package store

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
)

// Store manages the local Proven Skill Artifact store.
type Store struct {
	root string
}

// New creates a new Store. If storePath is empty, it defaults to
// the PSK_STORE environment variable or ~/.psk/store/.
func New(storePath string) *Store {
	if storePath == "" {
		storePath = os.Getenv("PSK_STORE")
	}
	if storePath == "" {
		home, _ := os.UserHomeDir()
		storePath = filepath.Join(home, ".psk", "store")
	}
	return &Store{root: storePath}
}

// Init creates the store directory and .store-version file if they don't exist.
func (s *Store) Init() error {
	if err := os.MkdirAll(s.root, 0o755); err != nil {
		return fmt.Errorf("failed to create store directory: %w", err)
	}
	versionFile := filepath.Join(s.root, ".store-version")
	if _, err := os.Stat(versionFile); os.IsNotExist(err) {
		if err := os.WriteFile(versionFile, []byte("1\n"), 0o644); err != nil {
			return fmt.Errorf("failed to write .store-version: %w", err)
		}
	}
	return nil
}

// Exists checks if a skill artifact with the given name and version already
// exists in the store.
func (s *Store) Exists(name, version string) bool {
	path := filepath.Join(s.root, name, version)
	_, err := os.Stat(path)
	return err == nil
}

// ArtifactPath returns the path where a skill artifact would be stored.
func (s *Store) ArtifactPath(name, version string) string {
	return filepath.Join(s.root, name, version)
}

// Add copies a skill directory into the store at {name}/{version}/.
// It writes atomically via a temp directory + os.Rename.
// If force is true, an existing artifact is replaced.
func (s *Store) Add(name, version, sourceDir string, manifest Manifest, force bool) (string, error) {
	destDir := filepath.Join(s.root, name, version)

	if !force {
		if _, err := os.Stat(destDir); err == nil {
			return "", fmt.Errorf("skill %s@%s already exists in store", name, version)
		}
	}

	// Create parent directory
	if err := os.MkdirAll(filepath.Join(s.root, name), 0o755); err != nil {
		return "", fmt.Errorf("failed to create skill directory: %w", err)
	}

	// Write to temp directory first for atomicity
	tmpDir := destDir + fmt.Sprintf(".tmp.%d", os.Getpid())
	if err := os.MkdirAll(tmpDir, 0o755); err != nil {
		return "", fmt.Errorf("failed to create temp directory: %w", err)
	}

	// Clean up temp dir on failure
	success := false
	defer func() {
		if !success {
			os.RemoveAll(tmpDir)
		}
	}()

	// Copy skill files from source to temp dir
	if err := copyDir(sourceDir, tmpDir); err != nil {
		return "", fmt.Errorf("failed to copy skill files: %w", err)
	}

	// Write manifest.json
	manifestPath := filepath.Join(tmpDir, "manifest.json")
	if err := WriteManifest(manifestPath, manifest); err != nil {
		return "", err
	}

	// If force, remove existing artifact
	if force {
		os.RemoveAll(destDir)
	}

	// Atomic rename
	if err := os.Rename(tmpDir, destDir); err != nil {
		return "", fmt.Errorf("failed to move artifact to store: %w", err)
	}

	success = true
	return destDir, nil
}

// List returns all manifests in the store, sorted by name then version.
func (s *Store) List() ([]Manifest, error) {
	var manifests []Manifest

	// Check if store exists
	if _, err := os.Stat(s.root); os.IsNotExist(err) {
		return manifests, nil
	}

	entries, err := os.ReadDir(s.root)
	if err != nil {
		return nil, fmt.Errorf("failed to read store: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		name := entry.Name()
		versionEntries, err := os.ReadDir(filepath.Join(s.root, name))
		if err != nil {
			continue
		}
		for _, ve := range versionEntries {
			if !ve.IsDir() {
				continue
			}
			manifestPath := filepath.Join(s.root, name, ve.Name(), "manifest.json")
			m, err := ReadManifest(manifestPath)
			if err != nil {
				continue
			}
			manifests = append(manifests, m)
		}
	}

	sort.Slice(manifests, func(i, j int) bool {
		if manifests[i].Name != manifests[j].Name {
			return manifests[i].Name < manifests[j].Name
		}
		return manifests[i].Version < manifests[j].Version
	})

	return manifests, nil
}

// copyDir copies all files and directories from src to dst.
func copyDir(src, dst string) error {
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			if err := os.MkdirAll(dstPath, 0o755); err != nil {
				return err
			}
			if err := copyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			if err := copyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}
	return nil
}

// copyFile copies a single file from src to dst.
func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return out.Close()
}
