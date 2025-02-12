package mango

import (
	"fmt"
	"net/http"
)

type GroupOption func(*Group)

type Group struct {
	app         *App
	middlewares []MiddlewareHandler
}

func WithMiddleware(middleware MiddlewareHandler) GroupOption {
	return func(g *Group) {
		g.middlewares = append(g.middlewares, middleware)
	}
}

func (g *Group) Use(middleware MiddlewareHandler) {
	g.middlewares = append(g.middlewares, middleware)
}

func (g *Group) GET(path string, handler func(c *Context) Response) {
	g.handle(http.MethodGet, path, handler)
}

func (g *Group) POST(path string, handler func(c *Context) Response) {
	g.handle(http.MethodPost, path, handler)
}

func (g *Group) PATCH(path string, handler func(c *Context) Response) {
	g.handle(http.MethodPatch, path, handler)
}

func (g *Group) PUT(path string, handler func(c *Context) Response) {
	g.handle(http.MethodPut, path, handler)
}

func (g *Group) DELETE(path string, handler func(c *Context) Response) {
	g.handle(http.MethodDelete, path, handler)
}

func (g *Group) ANY(path string, handler func(c *Context) Response) {
	g.handle("", path, handler)
}

func (g *Group) handle(method string, path string, handler func(c *Context) Response) {
	if method != "" {
		path = fmt.Sprintf("%s %s", method, path)
	}
	g.app.mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		rw := &responseWithStatus{w: w}
		c := &Context{
			req:         r,
			resp:        rw,
			app:         g.app,
			middlewares: g.middlewares,
		}
		rw.ctx = c
		c.handler = handler
		c.Next()
	})
}
