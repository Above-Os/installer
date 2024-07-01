package storage

// Table: install_config
const (
	queryFmtInsertInstallConfig = `
	  INSERT INTO %s (terminus_os_domainname, terminus_os_username, kube_type, vendor, gpu_enable, gpu_share, version)
		VALUES (?, ?, ?, ?, ?, ?, ?);`
	queryFmtInsertInstallLog = `
		INSERT INTO %s (message, state, percent)
		VALUES (?, ?, ?);`
	queryFmtQueryInstallState = `
		SELECT message, state, percent, created_at FROM %s ORDER BY id DESC LIMIT 1;
		`
)

const (
	querySQLiteSelectExistingTables = `
		SELECT name
		FROM sqlite_master
		WHERE type = 'table';`
)
