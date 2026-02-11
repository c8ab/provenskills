package skill

import (
	"fmt"
	"regexp"
	"strings"
)

var nameRegex = regexp.MustCompile(`^[a-z0-9]([a-z0-9-]*[a-z0-9])?$`)

// semverRegex matches major.minor.patch where each is a non-negative integer.
var semverRegex = regexp.MustCompile(`^\d+\.\d+\.\d+$`)

// majorMinorRegex matches major.minor where each is a non-negative integer.
var majorMinorRegex = regexp.MustCompile(`^\d+\.\d+$`)

// Validate checks a SkillFrontmatter against all validation rules.
// dirName is the name of the parent directory containing the SKILL.md.
// Returns a slice of human-readable error strings (empty if valid).
func Validate(fm SkillFrontmatter, dirName string) []string {
	var errs []string

	// Name validation
	if fm.Name == "" {
		errs = append(errs, "name: required field is missing")
	} else {
		if len(fm.Name) > 64 {
			errs = append(errs, fmt.Sprintf("name: must be 1-64 characters (got %d)", len(fm.Name)))
		}
		if strings.Contains(fm.Name, "--") {
			errs = append(errs, "name: contains consecutive hyphens")
		} else if !nameRegex.MatchString(fm.Name) {
			switch {
			case strings.ToLower(fm.Name) != fm.Name:
				errs = append(errs, "name: contains uppercase characters (must be lowercase)")
			case strings.HasPrefix(fm.Name, "-") || strings.HasSuffix(fm.Name, "-"):
				errs = append(errs, "name: must not start or end with a hyphen")
			default:
				errs = append(errs, fmt.Sprintf("name: %q does not match required pattern [a-z0-9][a-z0-9-]*[a-z0-9]", fm.Name))
			}
		}
		if fm.Name != dirName {
			errs = append(errs, fmt.Sprintf("name: %q does not match directory name %q", fm.Name, dirName))
		}
	}

	// Description validation
	if fm.Description == "" {
		errs = append(errs, "description: required field is missing")
	} else if len(fm.Description) > 1024 {
		errs = append(errs, fmt.Sprintf("description: must be at most 1024 characters (got %d)", len(fm.Description)))
	}

	// Version validation
	if fm.Metadata.Version == "" {
		errs = append(errs, "metadata.version: required field is missing")
	} else if !semverRegex.MatchString(fm.Metadata.Version) && !majorMinorRegex.MatchString(fm.Metadata.Version) {
		errs = append(errs, fmt.Sprintf("metadata.version: %q is not valid semver", fm.Metadata.Version))
	}

	// Author validation
	if fm.Metadata.Author == "" {
		errs = append(errs, "metadata.author: required field is missing")
	}

	return errs
}

// NormalizeVersion converts a version string to semver format.
// If the version is major.minor, it appends .0.
// If already semver, returns as-is.
func NormalizeVersion(version string) string {
	if majorMinorRegex.MatchString(version) {
		return version + ".0"
	}
	return version
}
