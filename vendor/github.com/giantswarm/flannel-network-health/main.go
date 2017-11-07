package main

import (
	"fmt"
	"github.com/giantswarm/flannel-network-health/flag"
	"github.com/giantswarm/flannel-network-health/server"
	"github.com/giantswarm/flannel-network-health/service"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/microkit/command"
	microserver "github.com/giantswarm/microkit/server"
	"github.com/giantswarm/microkit/transaction"
	"github.com/giantswarm/micrologger"
	"github.com/giantswarm/microstorage"
	"github.com/giantswarm/microstorage/memory"
	"github.com/spf13/viper"
	"os"
)

var (
	f           *flag.Flag = flag.New()
	description string     = "Flannel-network-health serves as health endpoint for network configuration created by flannel-operator."
	gitCommit   string     = "n/a"
	name        string     = "flannel-network-health"
	source      string     = "https://github.com/giantswarm/flannel-network-health"
)

func main() {
	err := mainWithError()
	if err != nil {
		panic(fmt.Sprintf("%#v\n", microerror.Mask(err)))
	}
}

func readEnv() error {
	// load conf from ENV
	f.Service.NetworkConfig.BridgeInterface = os.Getenv("NETWORK_BRIDGE_NAME")
	f.Service.NetworkConfig.FlannelInterface = os.Getenv("NETWORK_FLANNEL_DEVICE")
	f.Service.FlannelFile = os.Getenv("NETWORK_ENV_FILE_PATH")
	f.Service.ListenAddress = os.Getenv("LISTEN_ADDRESS")
	if f.Service.NetworkConfig.BridgeInterface == "" {
		return microerror.Maskf(invalidConfigError, "NETWORK_BRIDGE_NAME must not be empty")
	}
	if f.Service.NetworkConfig.FlannelInterface == "" {
		return microerror.Maskf(invalidConfigError, "NETWORK_FLANNEL_DEVICE must not be empty")
	}
	if f.Service.FlannelFile == "" {
		return microerror.Maskf(invalidConfigError, "NETWORK_ENV_FILE_PATH must not be empty")
	}
	if f.Service.ListenAddress == "" {
		return microerror.Maskf(invalidConfigError, "LISTEN_ADDRESS must not be empty")
	}
	return nil
}

func mainWithError() error {
	var err error

	// Create a new logger which is used by all packages.
	var newLogger micrologger.Logger
	{
		loggerConfig := micrologger.DefaultConfig()
		loggerConfig.IOWriter = os.Stdout
		newLogger, err = micrologger.New(loggerConfig)
		if err != nil {
			return err
		}
	}

	// We define a server factory to create the custom server once all command
	// line flags are parsed and all microservice configuration is storted out.
	newServerFactory := func(v *viper.Viper) microserver.Server {
		err = readEnv()
		if err != nil {
			panic(err)
		}
		// Create a new custom service which implements business logic.
		var newService *service.Service
		{
			serviceConfig := service.DefaultConfig()

			serviceConfig.Flag = f
			serviceConfig.Logger = newLogger

			serviceConfig.Description = description
			serviceConfig.GitCommit = gitCommit
			serviceConfig.Name = name
			serviceConfig.Source = source

			newService, err = service.New(serviceConfig)
			if err != nil {
				panic(err)
			}
		}

		var storage microstorage.Storage
		{
			storage, err = memory.New(memory.DefaultConfig())
			if err != nil {
				panic(err)
			}
		}

		var transactionResponder transaction.Responder
		{
			c := transaction.DefaultResponderConfig()
			c.Logger = newLogger
			c.Storage = storage

			transactionResponder, err = transaction.NewResponder(c)
			if err != nil {
				panic(err)
			}
		}

		// Create a new custom server which bundles our endpoints.
		var newServer microserver.Server
		{
			serverConfig := server.DefaultConfig()

			serverConfig.MicroServerConfig.Logger = newLogger
			serverConfig.MicroServerConfig.ServiceName = name
			serverConfig.MicroServerConfig.TransactionResponder = transactionResponder
			serverConfig.MicroServerConfig.Viper = v
			serverConfig.MicroServerConfig.ListenAddress = f.Service.ListenAddress
			serverConfig.Service = newService

			newServer, err = server.New(serverConfig)
			if err != nil {
				panic(err)
			}
		}

		return newServer
	}

	// Create a new microkit command which manages our custom microservice.
	var newCommand command.Command
	{
		commandConfig := command.DefaultConfig()

		commandConfig.Logger = newLogger
		commandConfig.ServerFactory = newServerFactory

		commandConfig.Description = description
		commandConfig.GitCommit = gitCommit
		commandConfig.Name = name
		commandConfig.Source = source

		newCommand, err = command.New(commandConfig)
		if err != nil {
			return microerror.Mask(err)
		}
	}

	newCommand.CobraCommand().Execute()

	return nil
}
