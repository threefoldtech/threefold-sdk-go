package config

import (
	"fmt"

	"github.com/threefoldtech/tfgrid-sdk-go/grid-compose/pkg"
)

type Config struct {
	Version  string                 `yaml:"version"`
	Networks map[string]pkg.Network `yaml:"networks"`
	Services map[string]pkg.Service `yaml:"services"`
	Storage  map[string]pkg.Storage `yaml:"storage"`
}

func NewConfig() *Config {
	return &Config{
		Networks: make(map[string]pkg.Network),
		Services: make(map[string]pkg.Service),
		Storage:  make(map[string]pkg.Storage),
	}
}

func (c *Config) ValidateConfig() (err error) {
	if c.Version == "" {
		return ErrVersionNotSet
	}

	for name, network := range c.Networks {
		if network.Type == "" {
			return fmt.Errorf("%w for network %s", ErrNetworkTypeNotSet, name)
		}
	}

	for name, service := range c.Services {
		if service.Flist == "" {
			return fmt.Errorf("%w for service %s", ErrServiceFlistNotSet, name)
		}

		if service.Resources.CPU == 0 {
			return fmt.Errorf("%w for service %s", ErrServiceCPUResourceNotSet, name)
		}

		if service.Resources.Memory == 0 {
			return fmt.Errorf("%w for service %s", ErrServiceMemoryResourceNotSet, name)
		}
	}

	for name, storage := range c.Storage {
		if storage.Type == "" {
			return fmt.Errorf("%w for storage %s", ErrStorageTypeNotSet, name)
		}

		if storage.Size == "" {
			return fmt.Errorf("%w for storage %s", ErrStorageSizeNotSet, name)
		}
	}

	return nil
}
