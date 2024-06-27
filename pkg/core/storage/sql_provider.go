package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"path"
	"time"

	"bytetrade.io/web3os/installer/pkg/core/logger"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
)

func NewSQLProvider(name, dataSourceName string) (provider SQLProvider) {
	dbName := path.Join(dataSourceName, providerDataSourceName)
	db, err := sqlx.Open(providerDriverName, fmt.Sprintf("file:%s?cache=shared", dbName))
	db.SetMaxIdleConns(1)
	provider = SQLProvider{
		db:         db,
		name:       name,
		driverName: providerDriverName,
		errOpen:    err,

		log: logger.GetLogger(),
	}

	return provider
}

type SQLProvider struct {
	db *sqlx.DB

	name       string
	driverName string
	schema     string
	errOpen    error

	log *zap.SugaredLogger

	// Utility.
	sqlSelectExistingTables string

	// Table: migrations.
	sqlInsertMigration       string
	sqlSelectMigrations      string
	sqlSelectLatestMigration string
}

func (p *SQLProvider) StartupCheck() (err error) {
	if p.errOpen != nil {
		return fmt.Errorf("error opening database: %w", p.errOpen)
	}

	for i := 0; i < 10; i++ {
		if err = p.db.Ping(); err == nil {
			break
		}

		time.Sleep(time.Millisecond * 500)
	}

	if err != nil {
		return fmt.Errorf("error pinging database: %w", err)
	}

	// + todo 这里测试下初始化表的操作
	ctx := context.Background()
	fmt.Println("---w1---")
	if err = p.SchemaMigrate(ctx, true); err != nil {
		return err
	}

	return nil
}

// BeginTX begins a transaction with the storage provider when applicable.
func (p *SQLProvider) BeginTX(ctx context.Context) (c context.Context, err error) {
	var tx *sql.Tx

	if tx, err = p.db.Begin(); err != nil {
		return nil, err
	}

	return context.WithValue(ctx, ctxKeyTransaction, tx), nil
}

// Commit performs a storage provider commit when applicable.
func (p *SQLProvider) Commit(ctx context.Context) (err error) {
	tx, ok := ctx.Value(ctxKeyTransaction).(*sql.Tx)

	if !ok {
		return errors.New("could not retrieve tx")
	}

	return tx.Commit()
}

// Rollback performs a storage provider rollback when applicable.
func (p *SQLProvider) Rollback(ctx context.Context) (err error) {
	tx, ok := ctx.Value(ctxKeyTransaction).(*sql.Tx)

	if !ok {
		return errors.New("could not retrieve tx")
	}

	return tx.Rollback()
}

// Close the underlying storage provider.
func (p *SQLProvider) Close() (err error) {
	return p.db.Close()
}
