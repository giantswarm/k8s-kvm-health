package service

import "github.com/giantswarm/flannel-network-health/flag/service/network"

type Service struct {
	NetworkConfig network.Network
	FlannelFile   string
	ListenAddress string
}
