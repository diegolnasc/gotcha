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

type ProviderWorker interface {
	Events(c *gin.Context)
}

func (s *ServerHandler) NewServer() {
	router := gin.Default()
	router.GET("/healthz", healthz)
	router.POST("/", events(s))
	router.Run(":3000")
}

func healthz(c *gin.Context) {
	c.JSON(200, gin.H{
		"Error":  false,
		"Status": "UP",
	})
}

func events(s *ServerHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		s.ProviderWorker.Events(c)
	}
}
