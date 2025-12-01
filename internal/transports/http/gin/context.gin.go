package gin

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GinContext struct {
	c *gin.Context
}

func (ctx GinContext) BindJSON(obj any) error   { return ctx.c.ShouldBindJSON(obj) }
func (ctx GinContext) BindURI(obj any) error    { return ctx.c.ShouldBindUri(obj) }
func (ctx GinContext) BindQuery(obj any) error  { return ctx.c.ShouldBindQuery(obj) }
func (ctx GinContext) BindHeader(obj any) error { return ctx.c.ShouldBindHeader(obj) }
func (ctx GinContext) Bind(obj any) error       { return ctx.c.ShouldBind(obj) }
func (ctx GinContext) JSON(code int, v any) error {
	ctx.c.JSON(code, v)
	return nil
}
func (ctx GinContext) Param(n string) string {
	return ctx.c.Param(n)
}
func (ctx GinContext) GetUserID() string {
	val, exists := ctx.c.Get("user_id")
	if !exists {
		return ""
	}
	if s, ok := val.(string); ok {
		return s
	}
	return ""
}
func (ctx GinContext) Get(key string) any {
	val, _ := ctx.c.Get(key)
	return val
}
func (ctx GinContext) Set(key string, value any) {
	ctx.c.Set(key, value)
}
func (ctx GinContext) GetContext() context.Context {
	return ctx.c.Request.Context()
}
func (ctx GinContext) GetHeader(key string) string {
	return ctx.c.GetHeader(key)
}
func (ctx GinContext) SetHeader(key, value string) {
	ctx.c.Header(key, value)
}
func (ctx GinContext) GetCookie(name string) (string, error) {
	cookie, err := ctx.c.Cookie(name)
	if err != nil {
		if err == http.ErrNoCookie {
			return "", nil
		}
		return "", err
	}
	return cookie, nil
}
func (ctx GinContext) SetCookie(name, value string, maxAge int, path, domain string, secure, httpOnly bool) {
	ctx.c.SetCookie(name, value, maxAge, path, domain, secure, httpOnly)
}
func (ctx GinContext) RemoveCookie(name string) {
	ctx.c.SetCookie(name, "", -1, "/", "", false, true)
}
func (ctx GinContext) GetClientIP() string {
	return ctx.c.ClientIP()
}
func (ctx GinContext) GetUserAgent() string {
	return ctx.c.Request.UserAgent()
}

func NewGinContext(c *gin.Context) GinContext {
	return GinContext{c: c}
}
