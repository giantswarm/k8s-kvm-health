package service

import (
	"github.com/giantswarm/flannel-network-health/flag"
	"github.com/giantswarm/flannel-network-health/flag/service/network"
	"github.com/giantswarm/microerror"
	"strings"
	"testing"
)

func Test_Flannel_ParseIP(t *testing.T) {
	tests := []struct {
		config             func(flannelFile []byte) (network.Network, error)
		flannelFileContent []byte
		expectedConfig     network.Network
		expectedErr        error
	}{
		// test 0
		{
			config: func(flannelFile []byte) (network.Network, error) {
				conf := DefaultConfig()
				conf.Flag = flag.New()
				err := conf.parseIPs(flannelFile)
				return conf.Flag.Service.NetworkConfig, err
			},
			expectedConfig: network.Network{FlannelIP: "172.23.3.64", BridgeIP: "172.23.3.65"},
			flannelFileContent: []byte(`FLANNEL_NETWORK=172.23.3.0/24
FLANNEL_SUBNET=172.23.3.65/30
FLANNEL_MTU=1450
FLANNEL_IPMASQ=false`),
			expectedErr: nil,
		},
		// test 1
		{
			config: func(flannelFile []byte) (network.Network, error) {
				conf := DefaultConfig()
				conf.Flag = flag.New()
				err := conf.parseIPs(flannelFile)
				return conf.Flag.Service.NetworkConfig, err
			},
			expectedConfig: network.Network{FlannelIP: "198.168.0.0", BridgeIP: "198.168.0.1"},
			flannelFileContent: []byte(`FLANNEL_NETWORK=198.168.0.0/24
FLANNEL_SUBNET=198.168.0.1/30
FLANNEL_MTU=1450
FLANNEL_IPMASQ=false`),
			expectedErr: nil,
		},
		// test 2 - missing FLANNEL_SUBNET
		{
			config: func(flannelFile []byte) (network.Network, error) {
				conf := DefaultConfig()
				conf.Flag = flag.New()
				err := conf.parseIPs(flannelFile)
				return conf.Flag.Service.NetworkConfig, err
			},
			expectedConfig: network.Network{},
			flannelFileContent: []byte(`FLANNEL_NETWORK=192.168.0.0/24
FLANNEL_MTU=1450
FLANNEL_IPMASQ=false`),
			expectedErr: invalidFlannelConfigurationError,
		},
		// test 3 - invalid subnet in flannel file
		{
			config: func(flannelFile []byte) (network.Network, error) {
				conf := DefaultConfig()
				conf.Flag = flag.New()
				err := conf.parseIPs(flannelFile)
				return conf.Flag.Service.NetworkConfig, err
			},
			expectedConfig: network.Network{},
			flannelFileContent: []byte(`FLANNEL_NETWORK=198.168.0.0/24
FLANNEL_SUBNET=_x.68.c.0/30
FLANNEL_MTU=1450
FLANNEL_IPMASQ=false`),
			expectedErr: invalidFlannelConfigurationError,
		},
		// test 4 - empty flannel file
		{
			config: func(flannelFile []byte) (network.Network, error) {
				conf := DefaultConfig()
				conf.Flag = flag.New()
				err := conf.parseIPs(flannelFile)
				return conf.Flag.Service.NetworkConfig, err
			},
			expectedConfig:     network.Network{},
			flannelFileContent: []byte(``),
			expectedErr:        invalidFlannelConfigurationError,
		},
		// test 5 - non flannel file
		{
			config: func(flannelFile []byte) (network.Network, error) {
				conf := DefaultConfig()
				conf.Flag = flag.New()
				err := conf.parseIPs(flannelFile)
				return conf.Flag.Service.NetworkConfig, err
			},
			expectedConfig: network.Network{},
			flannelFileContent: []byte(`machine:
  services:
    - docker

dependencies:
  override:
    - |
      wget -q $(curl -sS -H "Authorization: token $RELEASE_TOKEN" https://api.github.com/repos/giantswarm/architect/releases/latest | grep browser_download_url | head -n 1 | cut -d '"' -f 4)
    - chmod +x ./architect
    - ./architect version
`),
			expectedErr: invalidFlannelConfigurationError,
		},
	}

	for index, test := range tests {
		networkConfig, err := test.config(test.flannelFileContent)

		if microerror.Cause(err) != microerror.Cause(test.expectedErr) {
			t.Fatalf("%d: unexcepted error, expected %#v but got %#v", index, test.expectedErr, err)
		}
		if test.expectedErr == nil {
			if strings.Compare(networkConfig.FlannelIP, test.expectedConfig.FlannelIP) != 0 {
				t.Fatalf("%d: Incorrent ip, expected %s but got %s.", index, test.expectedConfig.FlannelIP, networkConfig.FlannelIP)
			}
			if strings.Compare(networkConfig.BridgeIP, test.expectedConfig.BridgeIP) != 0 {
				t.Fatalf("%d: Incorrent ip, expected %s but got %s.", index, test.expectedConfig.BridgeIP, networkConfig.BridgeIP)
			}
		}
	}
}
