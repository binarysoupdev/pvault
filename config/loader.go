package config

type Loader[Config any] interface {
	GetName() string
	Load() (Config, error)
}
