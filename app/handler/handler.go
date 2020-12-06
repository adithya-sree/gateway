package handler

import (
	"github.com/adithya-sree/gateway/app/file"
	"github.com/adithya-sree/gateway/config"
	"github.com/adithya-sree/logger"
)

// File Logger
var out = logger.GetLogger(config.LogFile, "handler")

// Request Handlers
type Handler struct {
	defs map[string]file.Definition
}

// Creates New Handler
func NewHandler() (*Handler, error) {
	// Read in Definitions
	defs, err := file.Read()
	if err != nil {
		return nil, err
	}
	// Get all definitions map
	enabled := getEnabledDefs(defs)
	// Create & return handler
	return &Handler{
		defs: enabled,
	}, nil
}

// Parses enabled definitions into map
func getEnabledDefs(defs []file.Definition) map[string]file.Definition {
	// Make & Allocate Memory for Definitions Map
	var enabled = make(map[string]file.Definition)
	// Loop through all read configs and add all enabled
	for _, d := range defs {
		if d.Enabled {
			enabled[d.ServiceRoute] = d
		}
	}
	// Return Enabled
	return enabled
}