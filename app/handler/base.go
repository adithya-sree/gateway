package handler

import (
	"github.com/adithya-sree/commons"
	"net/http"
	"time"
)

var startTime time.Time

// Set Start Time
func init() {
	startTime = time.Now()
}

// Uptime Response
type uptimeResponse struct {
	Start  string        `json:"start-time"`
	Uptime time.Duration `json:"uptime"`
}

// ECV Request
func (h Handler) Ecv() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		// Get Session
		sessionId := request.Context().Value("session").(string)
		out.Infof("[%s] - ECV Check Received", sessionId)
		// Return 200
		writer.WriteHeader(http.StatusOK)
	}
}

// Uptime Request
func (h *Handler) Uptime() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		sessionId := request.Context().Value("session").(string)
		out.Infof("[%s] - Uptime Request Received", sessionId)
		// Return 200
		_, _ = commons.RespondJSON(writer, http.StatusOK, uptimeResponse{
			Start:  startTime.Format("2006.01.02 15:04:05"),
			Uptime: time.Since(startTime),
		})
	}
}