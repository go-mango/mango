package mango

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

type AppOption func(*App)

func WithErrorHandler(handle func(*Context, error) Response) AppOption {
	return func(a *App) {
		a.handleError = handle
	}
}

type App struct {
	mux         *http.ServeMux
	addr        string
	handleError func(*Context, error) Response
	middlewares []MiddlewareHandler
}

func New(options ...AppOption) *App {
	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = ":8000"
	}
	a := &App{
		mux: http.NewServeMux(),
		handleError: func(c *Context, err error) Response {
			switch {
			case c.Status() >= 500:
				log.Println(err)
				return jsonRespond(c.Status(), map[string]any{"message": "internal server error"})
			default:
				return jsonRespond(c.Status(), map[string]any{"message": err.Error()})
			}
		},
		addr: addr,
	}
	for _, o := range options {
		o(a)
	}
	return a
}

func (a *App) Group(options ...GroupOption) *Group {
	mws := make([]MiddlewareHandler, len(a.middlewares))
	copy(mws, a.middlewares)
	g := &Group{
		app:         a,
		middlewares: mws,
	}
	for _, o := range options {
		o(g)
	}
	return g
}

func (a *App) Use(middleware MiddlewareHandler) {
	a.middlewares = append(a.middlewares, middleware)
}

func (a *App) GET(path string, handler func(c *Context) Response) {
	a.handle(http.MethodGet, path, handler)
}

func (a *App) POST(path string, handler func(c *Context) Response) {
	a.handle(http.MethodPost, path, handler)
}

func (a *App) PATCH(path string, handler func(c *Context) Response) {
	a.handle(http.MethodPatch, path, handler)
}

func (a *App) PUT(path string, handler func(c *Context) Response) {
	a.handle(http.MethodPut, path, handler)
}

func (a *App) DELETE(path string, handler func(c *Context) Response) {
	a.handle(http.MethodDelete, path, handler)
}

func (a *App) ANY(path string, handler func(c *Context) Response) {
	a.handle("", path, handler)
}

func (a *App) handle(method string, path string, handler func(c *Context) Response) {
	if method != "" {
		path = fmt.Sprintf("%s %s", method, path)
	}
	a.mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		rw := &responseWithStatus{w: w}
		c := &Context{
			req:         r,
			resp:        rw,
			app:         a,
			middlewares: a.middlewares,
		}
		rw.ctx = c
		c.handler = handler
		c.Next()
	})
}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.mux.ServeHTTP(w, r)
}

func (a *App) Listen() {
	panic(http.ListenAndServe(a.addr, a))
}
