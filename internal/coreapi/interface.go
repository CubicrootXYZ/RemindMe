package coreapi

import "github.com/gin-gonic/gin"

// CoreAPI provides an API for the core entities.
type CoreAPI interface {
	RegisterRoutes(*gin.Engine) error
}
