package storage

const (
	providerDriverName     = "sqlite3e"
	providerDataSourceName = "storage.db"
	providerSQLite         = "sqlite"
)

type ctxKey int

const (
	ctxKeyTransaction ctxKey = iota
)
