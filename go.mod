module github.com/giantswarm/k8s-kvm-health

go 1.14

require (
	github.com/giantswarm/microendpoint v0.0.0-20180904075734-f77c569259ae
	github.com/giantswarm/microerror v0.0.0-20181001144842-3bc3cb1a3670
	github.com/giantswarm/microkit v0.0.0-20181107110722-aaff79223ca0
	github.com/giantswarm/micrologger v0.0.0-20181005142804-22dc3b5565d1
	github.com/giantswarm/versionbundle v0.0.0-20181005143259-9a4f3249a5b5 // indirect
	github.com/go-kit/kit v0.10.0
	github.com/go-resty/resty v0.0.0-00010101000000-000000000000 // indirect
	github.com/juju/errgo v0.0.0-20140925100237-08cceb5d0b53 // indirect
	github.com/sparrc/go-ping v0.0.0-20181106165434-ef3ab45e41b0
	github.com/spf13/cast v1.3.0 // indirect
	github.com/spf13/pflag v1.0.3 // indirect
	github.com/spf13/viper v1.2.1
)

replace (
	github.com/go-resty/resty => gopkg.in/resty.v1 v1.11.0
	github.com/gogo/protobuf => github.com/gogo/protobuf v1.3.2
)
