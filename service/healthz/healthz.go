package healthz

import (
	"github.com/giantswarm/k8s-kvm-health/service/healthz/kvm"
	"github.com/giantswarm/microendpoint/service/healthz"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
)

// Config represents the configuration used to create a healthz service.
type Config struct {
	// Dependencies.
	CheckAPI  bool
	IPAddress string
	Logger    micrologger.Logger
}

// DefaultConfig provides a default configuration to create a new healthz
// service by best effort.
func DefaultConfig() Config {
	return Config{
		// Dependencies.
		IPAddress: "",
		Logger:    nil,
	}
}

// New creates a new configured healthz service.
func New(config Config) (*Service, error) {
	var err error

	var kvmService healthz.Service
	{
		var kvmServiceConfig kvm.Config
		kvmServiceConfig.IP = config.IPAddress
		kvmServiceConfig.Logger = config.Logger
		kvmServiceConfig.CheckAPI = config.CheckAPI
		kvmService, err = kvm.New(kvmServiceConfig)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	newService := &Service{
		KVM: kvmService,
	}

	return newService, nil
}

// Service is the healthz service collection.
type Service struct {
	KVM healthz.Service
}
