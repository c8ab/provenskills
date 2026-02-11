package cli

import (
	"fmt"
	"os"

	"github.com/c8ab/provenskills/internal/exitcode"
)

const version = "0.1.0"

const helpText = `psk - Package and manage Proven Skill Artifacts

Usage:
  psk <command> [flags]

Commands:
  build     Package a skill directory into an artifact
  list      List all skills in the local store
  validate  Validate a skill directory

Flags:
  --help      Show this help message
  --version   Show psk version

Environment:
  PSK_STORE   Override default store location (~/.psk/store/)`

// Run is the main entry point for the CLI. It parses the subcommand
// from args and dispatches to the appropriate handler.
// args should be os.Args (including the program name).
func Run(args []string) int {
	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, helpText)
		return exitcode.ErrGeneral
	}

	subcmd := args[1]

	switch subcmd {
	case "build":
		return RunBuild(args[2:])
	case "list":
		return RunList(args[2:])
	case "validate":
		return RunValidate(args[2:])
	case "--help", "-h", "help":
		fmt.Println(helpText)
		return exitcode.Success
	case "--version", "-v":
		fmt.Printf("psk version %s\n", version)
		return exitcode.Success
	default:
		fmt.Fprintf(os.Stderr, "error: unknown command %q\n\n%s\n", subcmd, helpText)
		return exitcode.ErrGeneral
	}
}
