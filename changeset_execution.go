package okapi

import "time"

// ChangesetExecution is an applied migration
type ChangesetExecution struct {
	ExecutionOrder int64
	ID             string
	Version        int64
	Checksum       string
	ExecutedAt     time.Time
	Status         string
	Comment        string
}

// SortableChangesetExecutions is needed to implement sort.Interface
type SortableChangesetExecutions []ChangesetExecution

func (m SortableChangesetExecutions) Len() int {
	return len(m)
}

func (m SortableChangesetExecutions) Swap(i int, j int) {
	m[i], m[j] = m[j], m[i]
}

func (m SortableChangesetExecutions) Less(i int, j int) bool {
	return m[i].ExecutionOrder < m[j].ExecutionOrder
}
