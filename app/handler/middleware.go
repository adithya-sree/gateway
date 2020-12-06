package handler

import (
	"context"
	"fmt"
	"github.com/adithya-sree/commons"
	"github.com/adithya-sree/gateway/app/file"
	"github.com/adithya-sree/gateway/config"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"net/http"
	"strings"
)

type RouteDTO struct {
	SessionId      string
	Definition     file.Definition
	QualifyingPath string
}

// Generates a Session ID
func (h Handler) SessionGenerator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get New UUID
		sessionId, err := uuid.NewUUID()
		if err != nil {
			// Return 500 if error generating session
			msg := fmt.Sprintf("An unexpected error has happened while trying to accquire a session id [%v]", err)
			out.Errorf("[CRITICAL] - [%s]", msg)
			_, _ = commons.RespondError(w, http.StatusInternalServerError, msg)
			return
		}
		out.Infof("[%s] - Created Session for Request", sessionId)
		// Add Session ID to Context
		ctx := context.WithValue(r.Context(), "session", sessionId.String())
		// Serve Next
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Authenticates a Request
func (h Handler) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Gets Session ID
		sessionId := r.Context().Value("session").(string)
		// Get Basic Auth from Request
		user, pass, _ := r.BasicAuth()
		if user != config.Username || pass != config.Password {
			// Return 401 if Username/Password does not match
			out.Infof("[%s] - Authentication failed for request", sessionId)
			_, _ = commons.RespondError(w, http.StatusUnauthorized, "missing/invalid credentials")
			return
		}
		// Serves next if authenticated
		out.Infof("[%s] - Authenticated request", sessionId)
		next.ServeHTTP(w, r)
	})
}

// Validates Definition
func (h *Handler) ValidateDefinition(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get Session
		sessionId := r.Context().Value("session").(string)
		// Get Definition
		def := chi.URLParam(r, "definition")
		// Validate Definition Exists
		definition := h.defs[def]
		if definition.ServiceName == "" {
			msg := fmt.Sprintf("[%s] - Unable to route request, no definitions found for [%s]", sessionId, def)
			out.Warnf(msg)
			_, _ = commons.RespondError(w, http.StatusNotFound, msg)
			return
		}
		// Serve Next with Context
		out.Infof("[%s] - Definition found for [%s] [%v]", sessionId, def, definition)
		ctx := context.WithValue(r.Context(), "session", &RouteDTO{
			SessionId:      sessionId,
			Definition:     definition,
			QualifyingPath: strings.Replace(r.URL.String(), "/gateway/"+definition.ServiceRoute, "", 1),
		})
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}