package storage

import (
	"context"

	"bytetrade.io/web3os/installer/pkg/model"
)

type Provider interface {
	StartupCheck() (err error)

	Ping() (err error)

	// Close the underlying storage provider.
	Close() (err error)

	SaveInstallConfig(ctx context.Context, config model.InstallModelReq) (err error)
}
