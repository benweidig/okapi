package okapi

type PostgresDialect struct {
}

func (d PostgresDialect) EnsureChangelog() string {
	return `
		CREATE TABLE IF NOT EXISTS okapi_changelog (
			execution_order		SERIAL PRIMARY KEY,
			id					VARCHAR(255)	NOT NULL,
			checksum			VARCHAR(32)		NOT NULL,
			comment				VARCHAR(255)	NOT NULL,
			executed_at			TIMESTAMP		NOT NULL,
			status				VARCHAR(16)		NOT NULL
		) Engine=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;`
}

func (d PostgresDialect) InsertRecord() string {
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

func (d PostgresDialect) ExecutedChangesets() string {
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
