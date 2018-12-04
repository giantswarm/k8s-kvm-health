package kvm

import (
	"context"
	"fmt"
	"github.com/giantswarm/microendpoint/service/healthz"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/sparrc/go-ping"
	"net/http"
	"time"
)

const (
	// Description describes which functionality this health check implements.
	Description = "Ensure KVM is responding to the assigned ip."
	// Name is the identifier of the health check. This can be used for emitting
	// metrics.
	Name = "kvmHealthz"

	// config
	pingCount = 1
)

// Config represents the configuration used to create a healthz service.
type Config struct {
	// Dependencies.
	CheckAPI bool
	IP       string
	Logger   micrologger.Logger
}

// Service implements the healthz service interface.
type Service struct {
	// Dependencies.
	checkAPI bool
	ip       string
	logger   micrologger.Logger

	// Settings.
	timeout time.Duration
}

// New creates a new configured healthz service.
func New(config Config) (*Service, error) {
	// Dependencies.
	if config.IP == "" {
		return nil, microerror.Maskf(invalidConfigError, "config.IP must not be empty string")
	}
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.Logger must not be empty")
	}

	newService := &Service{
		// Dependencies.
		checkAPI: config.CheckAPI,
		ip:       config.IP,
		logger:   config.Logger,
	}

	return newService, nil
}

// GetHealthz implements the health check for network interface.
func (s *Service) GetHealthz(ctx context.Context) (healthz.Response, error) {
	pingFailed, pingMsg := s.pingHealthCheck()

	response := healthz.Response{
		Description: Description,
		Failed:      pingFailed,
		Message:     pingMsg,
		Name:        Name,
	}

	// check api only if ping succeeded
	if !pingFailed && s.checkAPI {
		apiFailed, apiMsg := s.apiHealthCheck()
		response.Failed = apiFailed
		response.Message += apiMsg
	}

	return response, nil
}

// implementation fo the interface healthz logic for ping check
func (s *Service) pingHealthCheck() (bool, string) {
	var message string
	// ping kvm
	pinger, err := ping.NewPinger(s.ip)
	if err != nil {
		message = fmt.Sprintf("Failed to init pinger.")
		return true, message
	}
	// set fail values
	var failed = true
	message = fmt.Sprintf("Healthcheck for KVM has failed. KVM is not responding on  %s.", s.ip)

	pinger.Count = pingCount
	pinger.Timeout = time.Second * 1
	pinger.SetPrivileged(true)
	pinger.OnRecv = func(pkt *ping.Packet) {
		// we got positive response
		failed = false
		message = fmt.Sprintf("Healthcheck for KVM has been successful. KVM is live and responding. on %s.", s.ip)
	}

	pinger.Run()

	// exit
	return failed, message
}

// implementation of the interface healthz logic for k8s api check
func (s *Service) apiHealthCheck() (bool, string) {
	var message string
	url := fmt.Sprintf("https://%s", s.ip)
	// send request to k8s API
	_, err := http.Get(url)
	if err != nil {
		message = fmt.Sprintf("Failed to send http request to k8s API. %s", err)
		return true, message
	}

	message = fmt.Sprintf("Healthcheck for k8s API has been successful. K8s is live and responding. on %s.", url)
	// exit
	return true, message
}
