package main

import (
	"os"
	"syscall"

	"github.com/oxyno-zeta/golang-graphql-example/pkg/golang-graphql-example/authx/authentication"
	"github.com/oxyno-zeta/golang-graphql-example/pkg/golang-graphql-example/authx/authorization"
	"github.com/oxyno-zeta/golang-graphql-example/pkg/golang-graphql-example/business"
	"github.com/oxyno-zeta/golang-graphql-example/pkg/golang-graphql-example/config"
	"github.com/oxyno-zeta/golang-graphql-example/pkg/golang-graphql-example/database"
	"github.com/oxyno-zeta/golang-graphql-example/pkg/golang-graphql-example/email"
	"github.com/oxyno-zeta/golang-graphql-example/pkg/golang-graphql-example/lockdistributor"
	"github.com/oxyno-zeta/golang-graphql-example/pkg/golang-graphql-example/log"
	amqpbusmessage "github.com/oxyno-zeta/golang-graphql-example/pkg/golang-graphql-example/messagebus/amqp"
	"github.com/oxyno-zeta/golang-graphql-example/pkg/golang-graphql-example/metrics"
	"github.com/oxyno-zeta/golang-graphql-example/pkg/golang-graphql-example/signalhandler"
	"github.com/oxyno-zeta/golang-graphql-example/pkg/golang-graphql-example/tracing"
	"github.com/oxyno-zeta/golang-graphql-example/pkg/golang-graphql-example/version"
)

func setupExtraServices(targets []string, sv *services) {}

func setupBusinessServices(targets []string, sv *services) { //nolint: unparam // Ignored
	// Create business services
	busServices := business.NewServices(sv.logger, sv.db, sv.authorizationSvc, sv.ldSvc)
	// Save
	sv.busServices = busServices
}

func setupBasicsServices(targets []string, sv *services) { //nolint: unparam // Ignored
	// Create variables for mandatory services
	logger := sv.logger
	cfgManager := sv.cfgManager

	// Create metrics client
	metricsCl := metrics.NewMetricsClient()
	// Save
	sv.metricsCl = metricsCl

	// Generate tracing service instance
	tracingSvc, err := tracing.New(cfgManager, logger)
	// Check error
	if err != nil {
		logger.Fatal(err)
	}
	// Prepare on reload hook
	cfgManager.AddOnChangeHook(func() {
		err = tracingSvc.Reload()
		// Check error
		if err != nil {
			logger.Fatal(err)
		}
	})
	// Save
	sv.tracingSvc = tracingSvc

	// Create database service
	db := database.NewDatabase("main", cfgManager, logger, metricsCl, tracingSvc)
	// Connect to engine
	err = db.Connect()
	// Check error
	if err != nil {
		logger.Fatal(err)
	}
	// Add configuration reload hook
	cfgManager.AddOnChangeHook(func() {
		err = db.Reconnect()
		// Check error
		if err != nil {
			logger.Fatal(err)
		}
	})
	// Save
	sv.db = db

	// Create new mail service
	mailSvc := email.NewService(cfgManager, logger)
	// Try to connect
	err = mailSvc.Initialize()
	// Check error
	if err != nil {
		logger.Fatal(err)
	}
	// Add configuration reload hook
	cfgManager.AddOnChangeHook(func() {
		err = mailSvc.Initialize()
		// Check error
		if err != nil {
			logger.Fatal(err)
		}
	})
	// Save
	sv.mailSvc = mailSvc

	// Create lock distributor service
	ld := lockdistributor.NewService(cfgManager, db)
	// Initialize lock distributor
	err = ld.Initialize(logger)
	// Check error
	if err != nil {
		logger.Fatal(err)
	}
	// Add configuration reload hook
	cfgManager.AddOnChangeHook(func() {
		err = ld.Initialize(logger)
		// Check error
		if err != nil {
			logger.Fatal(err)
		}
	})
	// Save
	sv.ldSvc = ld

	// Create signal handler service
	signalHandlerSvc := signalhandler.NewClient(logger, true, []os.Signal{syscall.SIGTERM, syscall.SIGINT})
	// Initialize service
	err = signalHandlerSvc.Initialize()
	// Check error
	if err != nil {
		logger.Fatal(err)
	}
	// Register closing database connections on system stop
	signalHandlerSvc.OnExit(func() {
		err = db.Close()
		// Check error
		if err != nil {
			logger.Fatal(err)
		}
	})
	// Save
	sv.signalHandlerSvc = signalHandlerSvc

	// Get config
	cfg := cfgManager.GetConfig()
	// Initialize
	var amqpSvc amqpbusmessage.Client
	// Check if amqp have configuration set
	if cfg.AMQP != nil {
		// Create amqp bus message service
		amqpSvc = amqpbusmessage.New(logger, cfgManager, tracingSvc, signalHandlerSvc, metricsCl)
		// Connect
		err = amqpSvc.Connect()
		// Check error
		if err != nil {
			logger.Fatal(err)
		}
		// Add configuration reload hook
		cfgManager.AddOnChangeHook(func() {
			err = amqpSvc.Reconnect()
			// Check error
			if err != nil {
				logger.Fatal(err)
			}
		})
		// Register closing connections on system stop
		signalHandlerSvc.OnExit(func() {
			err2 := amqpSvc.Close()
			// Check error
			if err2 != nil {
				logger.Fatal(err2)
			}
		})
		// Register canceling consumers on SIGTERM or SIGINT
		signalHandlerSvc.OnSignal(syscall.SIGTERM, func() {
			err2 := amqpSvc.CancelAllConsumers()
			// Check error
			if err2 != nil {
				logger.Fatal(err2)
			}
		})
		signalHandlerSvc.OnSignal(syscall.SIGINT, func() {
			err2 := amqpSvc.CancelAllConsumers()
			// Check error
			if err2 != nil {
				logger.Fatal(err2)
			}
		})
	}
	// Save
	sv.amqpSvc = amqpSvc

	// Create authentication service
	authoSvc := authorization.NewService(cfgManager)
	// Save
	sv.authorizationSvc = authoSvc

	// Create authentication service
	authenticationSvc := authentication.NewService(cfgManager)
	// Save
	sv.authenticationSvc = authenticationSvc
}

func setupMandatoryServices(sv *services) {
	// Create new logger
	logger := log.NewLogger()
	// Save
	sv.logger = logger

	// Create configuration manager
	cfgManager := config.NewManager(logger)

	// Load configuration
	err := cfgManager.Load()
	// Check error
	if err != nil {
		logger.Fatal(err)
	}

	// Get configuration
	cfg := cfgManager.GetConfig()
	// Configure logger
	err = logger.Configure(cfg.Log.Level, cfg.Log.Format, cfg.Log.FilePath)
	// Check error
	if err != nil {
		logger.Fatal(err)
	}

	// Watch change for logger (special case)
	cfgManager.AddOnChangeHook(func() {
		// Get configuration
		newCfg := cfgManager.GetConfig()
		// Configure logger
		err = logger.Configure(newCfg.Log.Level, newCfg.Log.Format, newCfg.Log.FilePath)
		// Check error
		if err != nil {
			logger.Fatal(err)
		}
	})
	// Save
	sv.cfgManager = cfgManager

	// Getting version
	v := version.GetVersion()
	// Save
	sv.version = v
}