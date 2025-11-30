package http

import (
	"go-modular-monolith/internal/app/core"
	"go-modular-monolith/internal/domain/product"
	"go-modular-monolith/internal/domain/user"

	transportGin "go-modular-monolith/internal/transports/http/gin"

	"github.com/gin-gonic/gin"
)

func NewGinServer(c *core.Container) *gin.Engine {
	r := gin.Default()

	v1 := r.Group("/v1")
	routes := NewRoutes(
		c.ProductHandler,
		c.UserHandler,
	)
	for _, route := range *routes {
		switch route.Handler.(type) {
		case func(product.Context) error:
			transportGin.AdapterToGinRoutes(v1, &route, func(ctx *gin.Context) product.Context {
				return transportGin.NewGinContext(ctx)
			})
		case func(user.Context) error:
			transportGin.AdapterToGinRoutes(v1, &route, func(ctx *gin.Context) user.Context {
				return transportGin.NewGinContext(ctx)
			})
		}
	}
	return r
}
