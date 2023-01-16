package handler

import (
	"github.com/diegolnasc/gotcha/pkg/model"
	"github.com/gin-gonic/gin"
)

// Server command configuration.
type ServerHandler struct {
	ProviderWorker ProviderWorker
	Config         model.Settings
}

// ProviderWorker events received by the implemented provider.
type ProviderWorker interface {
	Events(c *gin.Context)
}

// NewServer starts the middleware route.
func (s *ServerHandler) NewServer() {
	router := gin.Default()
	router.GET("/healthz", healthz)
	router.POST("/", events(s))
	if err := router.Run(":3000"); err != nil {
		panic("Unable to start the gotcha server")
	}
}

// healthz default health check response.
func healthz(c *gin.Context) {
	c.JSON(200, gin.H{
		"Error":  false,
		"Status": "UP",
	})
}

// events middleware to provider
func events(s *ServerHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		s.ProviderWorker.Events(c)
	}
}
