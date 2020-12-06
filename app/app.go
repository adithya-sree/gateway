package app

import (
	"github.com/adithya-sree/gateway/app/handler"
	"github.com/adithya-sree/gateway/config"
	"github.com/adithya-sree/logger"
	"github.com/go-chi/chi"
	"net/http"
)

// File Logger
var out = logger.GetLogger(config.LogFile, "app")

// Application Struct
type App struct {
	Router  *chi.Mux
	handler *handler.Handler
}

// Run Application
func (a *App) Run() {
	// Start Refresh of Configs
	a.handler.StartRefresh()
	// Start Listening on Port
	out.Infof("Starting to listen on port [%s]", config.Port)
	err := http.ListenAndServe(":"+config.Port, a.Router)
	out.Errorf("Service has stopped listening [%v]", err)
}

// Creates New App
func NewApp() (*App, error) {
	// Create Chi Mux
	r := chi.NewMux()
	// Create Request Handler
	h, err := handler.NewHandler()
	if err != nil {
		return nil, err
	}
	// Initial Application Routes
	r.Group(func(r chi.Router) {
		// Session Generation Middleware
		r.Use(h.SessionGenerator)
		// Base Route
		r.Get("/", h.Ecv())
		// ECV Route
		r.Get("/ecv", h.Ecv())
		// Update Route
		r.Get("/uptime", h.Uptime())
		// Admin Requests
		r.Route("/admin", func(r chi.Router) {
			// Refresh Route
			r.Post("/refresh", h.Refresh())
			// Definitions Route
			r.Get("/definitions", h.Definitions())
		})
		// Gateway Requests
		r.Route("/gateway/{definition:[a-z-]+}", func(r chi.Router) {
			// Validate Definition Middleware
			r.Use(h.ValidateDefinition)
			// Gateway Routes
			r.Get("/*", h.GatewayRouterV1())
			r.Post("/*", h.GatewayRouterV1())
			r.Put("/*", h.GatewayRouterV1())
			r.Delete("/*", h.GatewayRouterV1())
		})
	})
	// Return Application with Mux
	return &App{
		Router:  r,
		handler: h,
	}, nil
}