package okapi

import "fmt"

// DuplicateError is thrown if an ID is not unique
type DuplicateError struct {
	ID string
}

func (e DuplicateError) Error() string {
	return fmt.Sprintf("multiple changesets have the same ID '%s'", e.ID)
}

// ChecksumError is thrown if there's a checksum mismatch between IDs
type ChecksumError struct {
	ID string
}

func (e ChecksumError) Error() string {
	return fmt.Sprintf("invalid checksum for changeset ID '%s'", e.ID)
}

// MissingError is thrown if a changeset is missing (e.g. really missing, or wrong order)
type MissingError struct {
	ID string
}

func (e MissingError) Error() string {
	return fmt.Sprintf("changeset is missing/wrong order ID '%s'", e.ID)
}

// DialectNotFoundError is thrown if a Dialect isn't found
type DialectNotFoundError struct {
	DriverName string
}

func (e DialectNotFoundError) Error() string {
	return fmt.Sprintf("dialect not found for driver '%s'", e.DriverName)
}

type InvalidStatusError struct {
	Status ExecutionStatus
}

func (e InvalidStatusError) Error() string {
	return fmt.Sprintf("last executed changeset ended in status %s", executionStatusMap[e.Status])
}

type InvalidChangesetError struct {
	Reason string
}

func (e InvalidChangesetError) Error() string {
	return fmt.Sprintf("invalid changeset: %s", e.Reason)
}
