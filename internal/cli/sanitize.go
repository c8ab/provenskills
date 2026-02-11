package cli

import (
	"path/filepath"
	"strings"
)

// containsPathTraversal checks if a path contains ".." segments,
// which could indicate a path traversal attack.
func containsPathTraversal(path string) bool {
	cleaned := filepath.Clean(path)
	parts := strings.Split(cleaned, string(filepath.Separator))
	for _, part := range parts {
		if part == ".." {
			return true
		}
	}
	return false
}
