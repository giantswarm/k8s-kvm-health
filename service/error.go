package service

import (
	"github.com/giantswarm/microerror"
)

var invalidConfigError = microerror.New("invalid config")

// IsInvalidConfig asserts invalidConfigError.
func IsInvalidConfig(err error) bool {
	return microerror.Cause(err) == invalidConfigError
}

var invalidFlannelFileError = microerror.New("invalid kvm file")

func IsInvalidFlannelFile(err error) bool {
	return microerror.Cause(err) == invalidFlannelFileError
}

var invalidFlannelConfigurationError = microerror.New("invalid kvm configuration")

func IsInvalidFlannelConfiguration(err error) bool {
	return microerror.Cause(err) == invalidFlannelConfigurationError
}

var failedParsingFlannelSubnetError = microerror.New("failed parsing kvm file")

func IsFailedParsingFlannelSubnet(err error) bool {
	return microerror.Cause(err) == failedParsingFlannelSubnetError
}
