package okapi

import (
	"database/sql"
	"sort"
	"sync"
	"time"
)

// Okapi is the brain of the operation
type Okapi struct {
	driver     Driver
	changesets []Changeset
	mtx        sync.RWMutex
}

// New creates a new okapi
func New(d Driver, c []Changeset) *Okapi {
	return &Okapi{
		driver:     d,
		changesets: c,
	}
}

// WithSQLDriver is a convenience method for created a generic sql-based okapi by driver name
func WithSQLDriver(db *sql.DB, driverName string, c []Changeset) (*Okapi, error) {
	driver, err := newSQLDriver(db, driverName)
	if err != nil {
		return nil, err
	}
	return &Okapi{
		driver:     driver,
		changesets: c,
	}, nil
}

// Migrate starts the migration process, including creation of the changelog and validation
func (o *Okapi) Migrate() error {
	o.mtx.Lock()
	defer o.mtx.Unlock()

	// Step 1: Make sure the changelog table exists
	notify("ensuring changelog table exists", nil, nil)
	err := o.driver.EnsureChangelog()
	if err != nil {
		notify("couldn't ensure changelog table existance", nil, err)
		return err
	}

	// Step 2: Retrieve the already applied migrations, we need them later
	notify("retrieving already executed changesets", nil, nil)
	applied, err := o.driver.ExecutedChangesets()
	if err != nil {
		notify("couldn't retrieve already executed changesets", nil, err)
		return err
	}
	sort.Sort(SortableChangesetExecutions(applied))

	// Step 3: Validate the current state
	notify("validating changelog/changesets", nil, nil)
	err = o.validate(applied)
	if err != nil {
		notify("validation error", nil, err)
		return err
	}

	// Step 4: Build the migration plan
	notify("building migration plan", nil, nil)
	plan := o.migrationPlan(applied)

	// Step 5: Execute the plan!
	notify("executing migration plan", nil, nil)
	for _, m := range plan {
		_, execErr := o.driver.Exec(m.Script)
		now := time.Now()
		status := ExecutionStatusExecuted
		if execErr != nil {
			if m.SkipOnError {
				status = ExecutionStatusSkipped
			} else {
				status = ExecutionStatusFailed
			}
		}
		r := &ChangesetExecution{
			ID:         m.ID,
			Comment:    m.Comment,
			Checksum:   m.checksum(),
			Status:     executionStatusMap[status],
			ExecutedAt: now,
		}
		err = o.driver.InsertRecord(r)
		if err != nil {
			notify("error executing", r, err)
			return err
		}
		if status == ExecutionStatusFailed {
			notify("error executing", r, err)
			return execErr
		}

		notify("executed changeset", r, nil)
	}

	notify("migration finished", nil, nil)

	return nil
}

func (o *Okapi) validate(executed []ChangesetExecution) error {
	// Are we in an undesired state, e.g. last migration left failed changeset?
	if len(executed) > 0 {
		last := executed[len(executed)-1]
		if last.Status == executionStatusMap[ExecutionStatusFailed] {
			return InvalidStatusError{ExecutionStatusFailed}
		}
	}

	// Validate changesets
	idMap := map[string]Changeset{}
	for _, c := range o.changesets {
		// Ask the driver if changeset is valid
		err := o.driver.Validate(c)
		if err != nil {
			return err
		}
		// Check for duplicate IDs
		_, ok := idMap[c.ID]
		if ok {
			return DuplicateError{c.ID}
		}
		idMap[c.ID] = c
	}

	// Validate already executed
	for _, c := range executed {
		current, ok := idMap[c.ID]
		// Check if a migration is missing/out of order
		if ok == false || current.ID != c.ID {
			return MissingError{c.ID}
		}
		// Validate checksum
		if current.checksum() != c.Checksum {
			return ChecksumError{current.ID}
		}
	}

	return nil
}

func (o *Okapi) migrationPlan(applied []ChangesetExecution) []Changeset {
	// If we never ran migrations before we can safely run all the changesets
	if len(applied) == 0 {
		return o.changesets
	}

	// If we already executed changesets we only wont the latest. We can just cut them
	// out of the slice because the validation ran beforehand.
	return o.changesets[len(applied):]
}
