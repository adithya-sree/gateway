package handler

import (
	"fmt"
	"github.com/adithya-sree/commons"
	"github.com/adithya-sree/gateway/app/file"
	"github.com/adithya-sree/gateway/config"
	"net/http"
	"time"
)

const refresh = "[REFRESH]"

// Refresh Configuration Route
func (h *Handler) Refresh() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		// Get Session
		sessionId := request.Context().Value("session").(string)
		out.Infof("[%s] - Refresh Configurations Request", sessionId)
		// Read in new definitions
		newDefinitions, err := file.Read()
		if err != nil {
			// Return 500 if error configurations
			msg := fmt.Sprintf("[%s] - Error reading config file [%v]", sessionId, err)
			out.Errorf(msg)
			_, _ = commons.RespondError(writer, http.StatusInternalServerError, msg)
			return
		}
		// Set New Configurations
		definitions := getEnabledDefs(newDefinitions)
		h.defs = definitions
		// Return 200
		msg := fmt.Sprintf("[%s] - Refreshed Configurations, [%d] Enabled Services", sessionId, len(definitions))
		out.Info(msg)
		_, _ = commons.RespondSuccess(writer, http.StatusOK, msg)
	}
}

// Definitions Routes
func (h *Handler) Definitions() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		// Get Session
		sessionId := request.Context().Value("session").(string)
		out.Infof("[%s] - Definitions Request", sessionId)
		// Return 200
		_, _ = commons.RespondJSON(writer, http.StatusOK, h.defs)
	}
}

// Starts Refresh Routine
func (h *Handler) StartRefresh() {
	// Update Ticker
	var update = time.NewTicker(time.Duration(config.RefreshRate) * time.Second)
	// Run Routine
	go func() {
		// For Ticker Range
		for range update.C {
			out.Infof("%s - Refreshing Gateway Configuration", refresh)
			// Read in New Definitions
			newDefinitions, err := file.Read()
			if err != nil {
				out.Errorf("%s - Error while refreshing configs, using previous configs [%v]", refresh, err)
				continue
			}
			// Set New Configurations
			definitions := getEnabledDefs(newDefinitions)
			h.defs = definitions
			out.Infof("%s - Refreshed Configurations, [%d] Enabled Services", refresh, len(definitions))
		}
	}()
}

