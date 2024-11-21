package mango

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

type App struct {
	mux         *http.ServeMux
	addr        string
	validate    func(any) error
	handleError func(*Context, error) Response
	middlewares []MiddlewareHandler
}

func New(options ...Option) *App {
	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = ":8000"
	}
	a := &App{
		mux:      http.NewServeMux(),
		validate: func(a any) error { return nil },
		handleError: func(c *Context, err error) Response {
			switch {
			case c.Status() >= 500:
				log.Println(err)
				return text(c.Status(), "internal server error")
			default:
				return text(c.Status(), err.Error())
			}
		},
		addr: addr,
	}
	for _, o := range options {
		o(a)
	}
	return a
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
			req:  r,
			resp: rw,
			app:  a,
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
