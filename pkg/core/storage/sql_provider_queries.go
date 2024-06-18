package storage

const (
	querySQLiteSelectExistingTables = `
		SELECT name
		FROM sqlite_master
		WHERE type = 'table';`
)
