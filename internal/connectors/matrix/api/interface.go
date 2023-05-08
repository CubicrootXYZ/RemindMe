package api

import "github.com/gin-gonic/gin"

// API provides an interface for the matrix connector API.
type API interface {
	RegisterRoutes(*gin.Engine) error
}
