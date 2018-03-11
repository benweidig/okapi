package okapi

type MysqlDialect struct {
}

func (d MysqlDialect) EnsureChangelog() string {
	return `
		CREATE TABLE IF NOT EXISTS okapi_changelog (
			execution_order		INT(11) 		NOT NULL AUTO_INCREMENT PRIMARY KEY,
			id					VARCHAR(255)	NOT NULL,
			checksum			VARCHAR(32)		NOT NULL,
			comment				VARCHAR(255)	NOT NULL,
			executed_at			DATETIME		NOT NULL,
			status				VARCHAR(16)		NOT NULL
		) Engine=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;`
}

func (d MysqlDialect) InsertRecord() string {
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

func (d MysqlDialect) ExecutedChangesets() string {
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
