package http

import (
	"go-modular-monolith/internal/domain/product"
	"go-modular-monolith/internal/domain/user"
	"go-modular-monolith/pkg/routes"
)

func NewRoutes(
	productHandler product.ProductHandler,
	userHandler user.UserHandler,
) *[]routes.Route {
	return &[]routes.Route{
		{
			Method:  "GET",
			Path:    "/product",
			Handler: productHandler.List,
		},
		{
			Method:  "POST",
			Path:    "/product",
			Handler: productHandler.Create,
		},

		// User CRUD
		{
			Method:  "GET",
			Path:    "/user",
			Handler: userHandler.List,
		},
		{
			Method:  "POST",
			Path:    "/user",
			Handler: userHandler.Create,
		},
		{
			Method:  "GET",
			Path:    "/user/:id",
			Handler: userHandler.Get,
		},
		{
			Method:  "PUT",
			Path:    "/user/:id",
			Handler: userHandler.Update,
		},
		{
			Method:  "DELETE",
			Path:    "/user/:id",
			Handler: userHandler.Delete,
		},
	}
}
