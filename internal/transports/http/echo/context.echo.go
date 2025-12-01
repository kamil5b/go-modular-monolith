package echo

import (
	"context"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

type EchoContext struct {
	c echo.Context
}

func (ctx EchoContext) BindJSON(obj any) error   { return ctx.c.Bind(obj) }
func (ctx EchoContext) BindURI(obj any) error    { return ctx.c.Bind(obj) }
func (ctx EchoContext) BindQuery(obj any) error  { return ctx.c.Bind(obj) }
func (ctx EchoContext) BindHeader(obj any) error { return ctx.c.Bind(obj) }
func (ctx EchoContext) Bind(obj any) error       { return ctx.c.Bind(obj) }
func (ctx EchoContext) JSON(code int, v any) error {
	return ctx.c.JSON(code, v)
}
func (ctx EchoContext) Param(n string) string {
	return ctx.c.Param(n)
}
func (ctx EchoContext) GetUserID() string {
	val := ctx.c.Get("user_id")
	if val == nil {
		return ""
	}
	if s, ok := val.(string); ok {
		return s
	}
	return ""
}
func (ctx EchoContext) Get(key string) any {
	return ctx.c.Get(key)
}
func (ctx EchoContext) Set(key string, value any) {
	ctx.c.Set(key, value)
}
func (ctx EchoContext) GetContext() context.Context {
	return ctx.c.Request().Context()
}
func (ctx EchoContext) GetHeader(key string) string {
	return ctx.c.Request().Header.Get(key)
}
func (ctx EchoContext) SetHeader(key, value string) {
	ctx.c.Response().Header().Set(key, value)
}
func (ctx EchoContext) GetCookie(name string) (string, error) {
	cookie, err := ctx.c.Cookie(name)
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			return "", nil
		}
		return "", err
	}
	return cookie.Value, nil
}
func (ctx EchoContext) SetCookie(name, value string, maxAge int, path, domain string, secure, httpOnly bool) {
	cookie := &http.Cookie{
		Name:     name,
		Value:    value,
		MaxAge:   maxAge,
		Path:     path,
		Domain:   domain,
		Secure:   secure,
		HttpOnly: httpOnly,
	}
	ctx.c.SetCookie(cookie)
}
func (ctx EchoContext) RemoveCookie(name string) {
	cookie := &http.Cookie{
		Name:   name,
		Value:  "",
		MaxAge: -1,
		Path:   "/",
	}
	ctx.c.SetCookie(cookie)
}
func (ctx EchoContext) GetClientIP() string {
	return ctx.c.RealIP()
}
func (ctx EchoContext) GetUserAgent() string {
	return ctx.c.Request().UserAgent()
}

func NewEchoContext(c echo.Context) EchoContext {
	return EchoContext{c: c}
}
