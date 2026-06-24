package config

type MockLoader[Config any] struct {
	Name   string
	Config Config
	Error  error
}

func (m MockLoader[Config]) GetName() string {
	return m.Name
}

func (m MockLoader[Config]) Load() (Config, error) {
	return m.Config, m.Error
}
