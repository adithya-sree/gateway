package main

import (
	"github.com/adithya-sree/gateway/app"
	"github.com/adithya-sree/gateway/config"
	"github.com/adithya-sree/logger"
	"sync"
)

// File Logger
var out = logger.GetLogger(config.LogFile, "main")

// Main
func main() {
	// Start Service
	out.Infof("Gateway Service is starting.")
	// Create WaitGroup
	wg := &sync.WaitGroup{}
	// Add Delta
	wg.Add(1)
	// Run Service on Routine
	go start(wg)
	// Wait Blocking
	wg.Wait()
}

// Runs Application Blocking
func start(wg *sync.WaitGroup) {
	// Defer Finishing Delta
	defer func() {
		out.Infof("Process is existing")
		wg.Done()
	}()
	// Create Application
	a, err := app.NewApp()
	if err != nil {
		out.Errorf("Unable to start application, error while initializing [%v]", err)
		return
	}
	// Run Application
	a.Run()
}