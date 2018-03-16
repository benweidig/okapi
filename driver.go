package okapi

import (
	"time"
)

// Driver is doing the actual work according to its dialect
type Driver interface {
	// Validate a changeset so it won't be a problem for a record (e.g. field sizes)
	Validate(c Changeset) error

	// EnsureChangelog creates the changelog table if necessary
	EnsureChangelog() error

	// LogExecution insert a new
	LogExecution(ex *ChangesetExecution) error

	// ExecutedChangesets returns the already executed changeset records
	ExecutedChangesets() ([]ChangesetExecution, error)

	// Exec runs an sql script
	Exec(sql string) (time.Duration, error)
}
