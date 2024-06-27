package storage

// Table: install_config
const (
	queryFmtInsertInstallConfig = `
	  INSERT INTO %s (terminus_os_domainname, terminus_os_username, kube_type, vendor, gpu_enable, gpu_share, version)
		VALUES (?, ?, ?, ?, ?, ?, ?);`
)

const (
	querySQLiteSelectExistingTables = `
		SELECT name
		FROM sqlite_master
		WHERE type = 'table';`
)
