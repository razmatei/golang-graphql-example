package main

import (
	"context"
	"flag"
	"strings"
	"sync"

	"emperror.dev/errors"
	"github.com/oxyno-zeta/golang-graphql-example/pkg/golang-graphql-example/authx/authentication"
	"github.com/oxyno-zeta/golang-graphql-example/pkg/golang-graphql-example/authx/authorization"
	"github.com/oxyno-zeta/golang-graphql-example/pkg/golang-graphql-example/business"
	"github.com/oxyno-zeta/golang-graphql-example/pkg/golang-graphql-example/config"
	"github.com/oxyno-zeta/golang-graphql-example/pkg/golang-graphql-example/database"
	"github.com/oxyno-zeta/golang-graphql-example/pkg/golang-graphql-example/email"
	lockdistributor "github.com/oxyno-zeta/golang-graphql-example/pkg/golang-graphql-example/lockdistributor/sql"
	"github.com/oxyno-zeta/golang-graphql-example/pkg/golang-graphql-example/log"
	amqpbusmessage "github.com/oxyno-zeta/golang-graphql-example/pkg/golang-graphql-example/messagebus/amqp"
	"github.com/oxyno-zeta/golang-graphql-example/pkg/golang-graphql-example/metrics"
	"github.com/oxyno-zeta/golang-graphql-example/pkg/golang-graphql-example/signalhandler"
	"github.com/oxyno-zeta/golang-graphql-example/pkg/golang-graphql-example/tracing"
	"github.com/oxyno-zeta/golang-graphql-example/pkg/golang-graphql-example/version"
	"github.com/samber/lo"
	"go.uber.org/automaxprocs/maxprocs"
)

type services struct {
	// Mandatory
	logger     log.Logger
	cfgManager config.Manager
	version    *version.AppVersion
	// Basics
	metricsSvc        metrics.Service
	tracingSvc        tracing.Service
	db                database.DB
	mailSvc           email.Service
	ldSvc             lockdistributor.Service
	signalHandlerSvc  signalhandler.Service
	amqpSvc           amqpbusmessage.Service
	authorizationSvc  authorization.Service
	authenticationSvc authentication.Service
	// Extra
	// Business
	busServices *business.Services
}

var targetDefinitionsMap = map[string]*targetDefinition{
	// Basics
	"migrate-db": migrateDBTarget,
	"server":     serverTarget,
	// Extra
}

// Those definitions are saving daemon definitions that will be launched with every target.
var daemonDefinitions = []*daemonDefinition{}

// WaitGroup is used to wait for the program to finish goroutines.
var (
	wg       sync.WaitGroup
	daemonWg sync.WaitGroup
)

func main() {
	// Compute possible targets
	possibleTargetValues := lo.Keys(targetDefinitionsMap)
	// Add "all" in those cases
	possibleTargetValues = append(possibleTargetValues, "all")

	// Init flags
	var targets arrayFlags
	// Init config folder path flag
	var configFolderPath string

	// Create target flag
	flag.Var(
		&targets,
		"target",
		"Represents the application target to be launched (possible values:"+strings.Join(
			possibleTargetValues,
			",",
		)+")",
	)
	flag.StringVar(
		&configFolderPath,
		"config-folder-path",
		config.DefaultMainConfigFolderPath,
		"Represents the configuration folder path",
	)
	// Parse flags
	flag.Parse()

	// Init services
	sv := &services{}

	// Setup mandatory services
	setupMandatoryServices(sv, configFolderPath)

	_, err := maxprocs.Set(maxprocs.Logger(sv.logger.Infof))
	// Check error
	if err != nil {
		sv.logger.Fatal(err)
	}

	// Catch any panic
	defer func() {
		// Catch panic
		if errI := recover(); errI != nil {
			// Panic caught => Log and exit
			// Try to cast error
			err, ok := errI.(error)
			// Check if cast wasn't ok
			if !ok {
				// Transform it
				err = errors.Errorf("%+v", errI)
			} else {
				// Map introduce stack trace
				err = errors.WithStack(err)
			}

			// Log
			sv.logger.Fatal(err)
		}
	}()

	// Defer sync
	defer sv.logger.Sync() //nolint: errcheck // This is part of the job

	sv.logger.Infof(
		"Application version: %s (git commit: %s) built on %s",
		sv.version.Version,
		sv.version.GitCommit,
		sv.version.BuildDate,
	)

	// Check if list is empty
	if len(targets) == 0 {
		// Add "all" for default values
		targets = append(targets, "all")
	}

	// Check if "all" is present with other things
	if lo.Contains(targets, "all") && len(targets) != 1 {
		// Reset to "all"
		targets = []string{"all"}
	}
	// Uniq targets
	targets = lo.Uniq(targets)
	// Check if target list have only accepted values
	for _, targetFlag := range targets {
		if !lo.Contains(possibleTargetValues, targetFlag) {
			sv.logger.Fatalf("target %s not supported", targetFlag)
		}
	}

	sv.logger.Infof("Starting application with targets: %s", targets)

	// Setup services
	setupBasicsServices(targets, sv)

	// Setup extra services
	setupExtraServices(targets, sv)

	// Setup business services
	setupBusinessServices(targets, sv)

	// Select targets and filter them by primary or not
	// Initialize target definitions lists
	primaryList := []*targetDefinition{}
	otherList := []*targetDefinition{}

	// Check if this is a "all" target
	if len(targets) == 1 && targets[0] == "all" {
		// Loop over all possible targets
		for _, tDef := range targetDefinitionsMap {
			// Check if acceptable in "all"
			if tDef.InAllTarget {
				// Check if primary
				if tDef.Primary {
					primaryList = append(primaryList, tDef)
				} else {
					otherList = append(otherList, tDef)
				}
			}
		}
	} else {
		// Loop over targets
		for _, target := range targets {
			// Get target definition
			tDef := targetDefinitionsMap[target]
			// Check if primary
			if tDef.Primary {
				primaryList = append(primaryList, tDef)
			} else {
				otherList = append(otherList, tDef)
			}
		}
	}

	// Add count for daemon wait group
	daemonWg.Add(len(daemonDefinitions))

	// Create cancellable daemon context
	dCtx, dCancel := context.WithCancel(context.TODO())

	for _, dDef := range daemonDefinitions {
		go func(dDef *daemonDefinition) {
			// Inform routine is completed
			defer daemonWg.Done()

			// Run target
			dDef.Run(dCtx, targets, sv)
		}(dDef)
	}

	// Start all primary targets
	for _, tDef := range primaryList {
		// Run
		tDef.Run(targets, sv)
	}

	// Add count of other targets for waiting group
	wg.Add(len(otherList))

	// Start all other targets
	for _, tDef := range otherList {
		// Start routine
		go func(tDef *targetDefinition) {
			// Inform routine is completed
			defer wg.Done()

			// Run target
			tDef.Run(targets, sv)
		}(tDef)
	}

	// Wait
	wg.Wait()

	// Cancel daemon context
	dCancel()

	// Wait all daemons
	daemonWg.Wait()
}
