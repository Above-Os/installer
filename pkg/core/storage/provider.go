package storage

type Provider interface {
	StartupCheck() (err error)
}
