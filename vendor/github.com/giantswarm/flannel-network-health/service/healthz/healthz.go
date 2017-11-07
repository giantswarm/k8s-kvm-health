package healthz

import (
	"github.com/giantswarm/flannel-network-health/flag/service/network"
	"github.com/giantswarm/flannel-network-health/service/healthz/flannel"
	"github.com/giantswarm/microendpoint/service/healthz"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
)

// Config represents the configuration used to create a healthz service.
type Config struct {
	// Dependencies.
	NetworkConfig network.Network
	Logger        micrologger.Logger
}

// DefaultConfig provides a default configuration to create a new healthz
// service by best effort.
func DefaultConfig() Config {
	return Config{
		// Dependencies.
		NetworkConfig: network.Network{},
		Logger:        nil,
	}
}

// New creates a new configured healthz service.
func New(config Config) (*Service, error) {
	var err error

	var bridgeService healthz.Service
	{
		bridgeServiceConfig := flannel.DefaultConfig()
		bridgeServiceConfig.Name = config.NetworkConfig.BridgeInterface
		bridgeServiceConfig.IP = config.NetworkConfig.BridgeIP
		bridgeServiceConfig.Logger = config.Logger
		bridgeService, err = flannel.New(bridgeServiceConfig)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}
	var flannelService healthz.Service
	{
		flannelServiceConfig := flannel.DefaultConfig()
		flannelServiceConfig.Name = config.NetworkConfig.FlannelInterface
		flannelServiceConfig.IP = config.NetworkConfig.FlannelIP
		flannelServiceConfig.Logger = config.Logger
		flannelService, err = flannel.New(flannelServiceConfig)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	newService := &Service{
		Bridge:  bridgeService,
		Flannel: flannelService,
	}

	return newService, nil
}

// Service is the healthz service collection.
type Service struct {
	Bridge  healthz.Service
	Flannel healthz.Service
}
