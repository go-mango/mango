package mango

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime"
	"strings"
	"sync"
)

type Context struct {
	req     *http.Request
	resp    http.ResponseWriter
	status  int
	app     *App
	handler func(*Context) Response
	cursor  int
	aborted bool
	store   map[string]any
	mu      sync.RWMutex
}

func (c *Context) Set(key string, value any) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.store == nil {
		c.store = make(map[string]any)
	}
	c.store[key] = value
}

func (c *Context) Get(key string) any {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.store[key]
}

func (c *Context) Status() int {
	return c.status
}

func (c *Context) Request() *http.Request {
	return c.req
}

func (c *Context) ResponseWriter() http.ResponseWriter {
	return c.resp
}

func (c *Context) Read(p []byte) (int, error) {
	return c.req.Body.Read(p)
}

func (c *Context) Write(p []byte) (int, error) {
	return c.resp.Write(p)
}

func (c *Context) Next() {
	if c.aborted {
		return
	}
	defer func() {
		if err := recover(); err != nil {
			c.aborted = true
			aerr, ok := err.(*abortError)
			if ok {
				c.status = aerr.status
				if aerr.err != nil {
					c.app.handleError(c, aerr.err).Send(c.ResponseWriter())
					return
				}
				c.ResponseWriter().WriteHeader(aerr.status)
				return
			}
			c.status = http.StatusInternalServerError
			buf := make([]byte, 1024)
			n := runtime.Stack(buf, false)
			lines := bytes.Split(buf[:n], []byte("\n"))
			var location string
			if len(lines) > 6 {
				location = strings.Split(string(bytes.TrimSpace(lines[6])), " ")[0]
			}
			c.app.handleError(c, fmt.Errorf("%v (%s)", err, location)).Send(c.ResponseWriter())
		}
	}()
	if c.cursor == len(c.app.middlewares) {
		c.handler(c).Send(c.ResponseWriter())
		return
	}
	cursor := c.cursor
	c.cursor++
	c.app.middlewares[cursor](c)
}
