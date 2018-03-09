package okapi

// Dialect provides the SQL needed for the differenct statements
type Dialect interface {

	// SQL for creating the changelog table
	EnsureChangelog() string

	// SQL for onserting a record, including the placeholders
	InsertRecord() string

	// SQL for selecting all already executed records
	ExecutedChangesets() string
}
