package okapi

// ExecutionStatus is an enum type for simpler usage
type ExecutionStatus int

const (
	// ExecutionStatusUnknown is the default value
	ExecutionStatusUnknown ExecutionStatus = iota

	// ExecutionStatusExecuted was successfully executed
	ExecutionStatusExecuted

	// ExecutionStatusFailed indicates a failed changeset
	ExecutionStatusFailed

	// ExecutionStatusSkipped indicates a skipped changeset after an error
	ExecutionStatusSkipped
)

var executionStatusMap = map[ExecutionStatus]string{
	ExecutionStatusUnknown:  "UNKNOWN",
	ExecutionStatusExecuted: "EXECUTED",
	ExecutionStatusFailed:   "FAILED",
	ExecutionStatusSkipped:  "SKIPPED",
}
