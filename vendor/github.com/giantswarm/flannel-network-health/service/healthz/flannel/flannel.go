package flannel

import (
	"context"
	"fmt"
	"github.com/giantswarm/flannel-network-health/service/healthz/flannel/key"
	"github.com/giantswarm/microendpoint/service/healthz"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/vishvananda/netlink"
	"time"
)

const (
	// Description describes which functionality this health check implements.
	Description = "Ensure network interface is present and has proper network configuration."
	// Name is the identifier of the health check. This can be used for emitting
	// metrics.
	Name = "interfaceHealthz"
)

// Config represents the configuration used to create a healthz service.
type Config struct {
	// Dependencies.
	Name   string
	IP     string
	Logger micrologger.Logger
}

// DefaultConfig provides a default configuration to create a new healthz service
// by best effort.
func DefaultConfig() Config {
	return Config{
		// Dependencies.
		Name:   "",
		IP:     "",
		Logger: nil,
	}
}

// Service implements the healthz service interface.
type Service struct {
	// Dependencies.
	name   string
	ip     string
	logger micrologger.Logger

	// Settings.
	timeout time.Duration
}

// New creates a new configured healthz service.
func New(config Config) (*Service, error) {
	// Dependencies.
	if config.IP == "" {
		return nil, microerror.Maskf(invalidConfigError, "config.NetworkInterface.IP must not be empty string")
	}
	if config.Name == "" {
		return nil, microerror.Maskf(invalidConfigError, "config.NetworkInterface.Name must not be empty string")
	}
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.Logger must not be empty")
	}

	newService := &Service{
		// Dependencies.
		name:   config.Name,
		ip:     config.IP,
		logger: config.Logger,
	}

	return newService, nil
}

// GetHealthz implements the health check for network interface.
func (s *Service) GetHealthz(ctx context.Context) (healthz.Response, error) {
	message := fmt.Sprintf("Healthcheck for interface %s has been successful. Interface is present and configured with ip %s.", s.name, s.ip)

	failed, message := s.healthCheck(message)

	response := healthz.Response{
		Description: Description,
		Failed:      failed,
		Message:     message,
		Name:        Name + " " + s.name,
	}

	return response, nil
}

// implementation fo the interface healthz logic
func (s *Service) healthCheck(message string) (bool, string) {
	// load interface
	bridge, err := netlink.LinkByName(s.name)
	if err != nil {
		message = fmt.Sprintf("Cant find interface %s. %s", s.name, err)
		return true, message
	}
	// check ip on interface
	ipList, err := netlink.AddrList(bridge, netlink.FAMILY_V4)
	if err != nil || len(ipList) == 0 {
		message = fmt.Sprintf("Missing ip %s on the interface %s.", s.ip, s.name)
		return true, message
	}
	// compare ip on interface
	if len(ipList) > 0 && key.GetInterfaceIP(ipList) != s.ip {
		message = fmt.Sprintf("Wrong ip on interface %s. Expected %s, but found %s.", s.name, s.ip, key.GetInterfaceIP(ipList))
		return true, message
	}
	// all good
	return false, message
}
