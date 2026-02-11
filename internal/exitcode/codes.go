// Package exitcode defines process exit codes for the psk CLI.
package exitcode

const (
	// Success indicates the operation completed successfully.
	Success = 0
	// ErrGeneral indicates an unexpected or internal error.
	ErrGeneral = 1
	// ErrValidation indicates a metadata validation failure.
	ErrValidation = 2
	// ErrConflict indicates a store conflict (name+version already exists).
	ErrConflict = 3
	// ErrIO indicates a filesystem or IO error.
	ErrIO = 4
)
