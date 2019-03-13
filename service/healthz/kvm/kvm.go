package kvm

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/giantswarm/microendpoint/service/healthz"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/sparrc/go-ping"
	"net/http"
	"net/url"
	"time"
)

const (
	// Description describes which functionality this health check implements.
	Description = "Ensure KVM is responding to the assigned ip."
	// Name is the identifier of the health check. This can be used for emitting
	// metrics.
	Name = "kvmHealthz"

	// config
	pingCount         = 1
	httpsScheme       = "https"
	httpScheme        = "http"
	k8sAPIPort        = 443
	k8sKubeletPort    = 10248
	maxIdleConnection = 10
	maxTimeout        = 4
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
	client   *http.Client
	ip       string
	logger   micrologger.Logger
	tr       *http.Transport

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

	tr := &http.Transport{
		MaxConnsPerHost:   maxIdleConnection,
		TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
		DisableKeepAlives: true,
		IdleConnTimeout:   maxTimeout,
	}

	client := &http.Client{
		Transport: tr,
		Timeout:   maxTimeout,
	}

	newService := &Service{
		// Dependencies.
		checkAPI: config.CheckAPI,
		client:   client,
		ip:       config.IP,
		logger:   config.Logger,
		tr:       tr,
	}

	return newService, nil
}

// GetHealthz Provides Healthz implementation to check health status of network
// interface. It performs following checks in given order:
//  - Ping configured IP.
//  - Check that Kubelet instance in configured IP responds to HTTP request.
//  - Check that K8s API in configured IP responds to HTTPS request.
func (s *Service) GetHealthz(ctx context.Context) (healthz.Response, error) {
	var apiFailed, kubeletFailed, pingFailed bool
	var apiMsg, kubeletMsg, pingMsg string
	pingFailed, pingMsg = s.pingHealthCheck()

	response := healthz.Response{
		Description: Description,
		Failed:      pingFailed,
		Message:     pingMsg,
		Name:        Name,
	}

	// check kubelet only if ping succeeded
	if !pingFailed {
		kubeletFailed, kubeletMsg = s.httpHealthCheck(k8sKubeletPort, httpScheme)
		response.Failed = kubeletFailed
		response.Message = kubeletMsg
	}

	// check api only if ping and kubelet succeeded
	if !pingFailed && !kubeletFailed && s.checkAPI {
		apiFailed, apiMsg = s.httpHealthCheck(k8sAPIPort, httpsScheme)
		response.Failed = apiFailed
		response.Message = apiMsg
	}

	return response, nil
}

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

func (s *Service) httpHealthCheck(port int, scheme string) (bool, string) {
	var message string
	u := url.URL{
		Host:   fmt.Sprintf("%s:%d", s.ip, port),
		Path:   "healthz",
		Scheme: scheme,
	}
	// be sure to close idle connection after health check is finished
	defer s.tr.CloseIdleConnections()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		panic(fmt.Sprintf("unable to construct health check request: %q", err))
	}

	// close connection after health check request (the TCP connection gets
	// closed by deferred s.tr.CloseIdleConnections()).
	req.Header.Add("Connection", "close")

	// send request to http endpoint
	_, err = s.client.Do(req)
	if err != nil {
		message = fmt.Sprintf("Failed to send http request to endpoint %s. %s", u.String(), err)
		return true, message
	}

	message = fmt.Sprintf("Healthcheck for http endpoint %s has been successful.", u.String())
	// exit
	return false, message
}
