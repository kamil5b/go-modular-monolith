package echo

import (
	"go-modular-monolith/pkg/routes"

	"github.com/labstack/echo/v4"
)

func AppRoutesToEchoRoutes[T any](
	e *echo.Group,
	routes *routes.Route,
	domainContext func(echo.Context) T,
) *echo.Group {
	switch routes.Method {
	case "GET":
		e.GET(routes.Path, func(ctx echo.Context) error {
			return routes.Handler.(func(T) error)(domainContext(ctx))
		})
	case "POST":
		e.POST(routes.Path, func(ctx echo.Context) error {
			return routes.Handler.(func(T) error)(domainContext(ctx))
		})
	case "PUT":
		e.PUT(routes.Path, func(ctx echo.Context) error {
			return routes.Handler.(func(T) error)(domainContext(ctx))
		})
	case "PATCH":
		e.PATCH(routes.Path, func(ctx echo.Context) error {
			return routes.Handler.(func(T) error)(domainContext(ctx))
		})
	case "DELETE":
		e.DELETE(routes.Path, func(ctx echo.Context) error {
			return routes.Handler.(func(T) error)(domainContext(ctx))
		})
	}
	return e
}
