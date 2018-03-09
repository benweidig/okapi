package okapi

type SqliteDialect struct {
}

func (d SqliteDialect) EnsureChangelog() string {
	return `
		CREATE TABLE IF NOT EXISTS okapi_changelog (
			execution_order		INTEGER 		PRIMARY KEY AUTOINCREMENT,
			id					VARCHAR(255)	NOT NULL,
			checksum			VARCHAR(32)		NOT NULL,
			comment				VARCHAR(255)	NOT NULL,
			executed_at			DATETIME		NOT NULL,
			status				VARCHAR(16)		NOT NULL,
			UNIQUE(id)
		);`
}

func (d SqliteDialect) InsertRecord() string {
	return `
		INSERT INTO okapi_changelog (
			id,
			checksum,
			comment,
			executed_at,
			status
		)
		VALUES (?, ?, ?, ?, ?);`
}

func (d SqliteDialect) ExecutedChangesets() string {
	return `
		SELECT
			execution_order,
			id,
			checksum,
			comment,
			executed_at,
			status
		FROM okapi_changelog
		ORDER BY execution_order ASC;`
}
