package service

import (
	"github.com/giantswarm/microerror"
)

var invalidConfigError = microerror.New("invalid config")

// IsInvalidConfig asserts invalidConfigError.
func IsInvalidConfig(err error) bool {
	return microerror.Cause(err) == invalidConfigError
}

var invalidFlannelFileError = microerror.New("invalid flannel file")

func IsInvalidFlannelFile(err error) bool {
	return microerror.Cause(err) == invalidFlannelFileError
}

var invalidFlannelConfigurationError = microerror.New("invalid flannel configuration")

func IsInvalidFlannelConfiguration(err error) bool {
	return microerror.Cause(err) == invalidFlannelConfigurationError
}

var failedParsingFlannelSubnetError = microerror.New("failed parsing flannel file")

func IsFailedParsingFlannelSubnet(err error) bool {
	return microerror.Cause(err) == failedParsingFlannelSubnetError
}
