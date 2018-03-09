package okapi

import (
	"database/sql"
	"fmt"
	"time"
)

// sqlDriver is a simple driver supporting stdlib sql
type sqlDriver struct {
	Db      *sql.DB
	Dialect Dialect
}

func newSQLDriver(db *sql.DB, driverName string) (*sqlDriver, error) {
	dialect, err := dialect(driverName)
	if err != nil {
		return nil, err
	}
	return &sqlDriver{
		Db:      db,
		Dialect: dialect,
	}, nil
}

func (d *sqlDriver) Validate(c Changeset) error {
	lID := len(c.ID)
	if lID == 0 || lID > 255 {
		return InvalidChangesetError{fmt.Sprintf("invalid ID length - valid: 1-255 - actual: %d", lID)}
	}

	if len(c.Script) == 0 {
		return InvalidChangesetError{"empty script"}
	}

	lComment := len(c.Comment)
	if lComment == 0 || lComment > 255 {
		return InvalidChangesetError{fmt.Sprintf("invalid Comment length - valid: 1-255 - actual: %d", lComment)}
	}

	return nil
}

func (d *sqlDriver) EnsureChangelog() error {
	err := d.tx(func(tx *sql.Tx) error {
		_, err := tx.Exec(d.Dialect.EnsureChangelog())
		return err
	})
	return err
}

func (d *sqlDriver) InsertRecord(r *ChangesetExecution) error {
	err := d.tx(func(tx *sql.Tx) error {
		_, err := tx.Exec(
			d.Dialect.InsertRecord(),
			r.ID,
			r.Checksum,
			r.Comment,
			r.ExecutedAt,
			r.Status,
		)
		return err
	})
	return err
}

func (d *sqlDriver) ExecutedChangesets() ([]ChangesetExecution, error) {
	var applied []ChangesetExecution

	rows, err := d.Db.Query(d.Dialect.ExecutedChangesets())
	if err != nil {
		return applied, err
	}

	for rows.Next() {
		r := ChangesetExecution{}

		rows.Scan(
			&r.ExecutionOrder,
			&r.ID,
			&r.Checksum,
			&r.Comment,
			&r.ExecutedAt,
			&r.Status)

		applied = append(applied, r)
	}
	return applied, nil
}

func (d *sqlDriver) Exec(sqlScript string) (time.Duration, error) {
	start := time.Now()

	err := d.tx(func(tx *sql.Tx) error {
		_, err := tx.Exec(sqlScript)
		return err
	})
	elapsed := time.Since(start)
	return elapsed, err
}

func (d *sqlDriver) tx(sqlFunc func(*sql.Tx) error) error {
	tx, err := d.Db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		p := recover()
		if p != nil {
			switch p := p.(type) {
			case error:
				err = p
			default:
				err = fmt.Errorf("%s", p)
			}
		}

		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	return sqlFunc(tx)
}

func dialect(driverName string) (Dialect, error) {
	switch driverName {
	case "sqlite3":
		return SqliteDialect{}, nil
	case "mysql":
		return MysqlDialect{}, nil
	default:
		return nil, DialectNotFoundError{DriverName: driverName}
	}
}
