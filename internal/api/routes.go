package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (server *server) assembleRoutes() error {
	router := gin.New()

	for name, provider := range server.config.EndpointProviders {
		server.logger.Infof("registering routes from: %s", name)

		err := provider.RegisterRoutes(router)
		if err != nil {
			server.logger.Errorf("error while setting up routes for '%s', routes might not be setup correctly: %v", name, err)
			continue
		}
	}

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	server.server = &http.Server{
		Addr:         server.config.Address,
		Handler:      router,
		ReadTimeout:  time.Second * 15,
		WriteTimeout: time.Second * 15,
	}

	return nil
}
