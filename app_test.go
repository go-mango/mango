package mango_test

import (
	"bytes"
	stdjson "encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-mango/json"
	"github.com/go-mango/mango"
	"github.com/go-mango/validate"
	"github.com/stretchr/testify/assert"
)

func run(t *testing.T, app *mango.App, method string, path string, body any, then func(*httptest.ResponseRecorder)) {
	var reader io.Reader
	if body != nil {
		b, _ := stdjson.Marshal(body)
		reader = bytes.NewReader(b)
	}
	req, err := http.NewRequest(method, path, reader)
	assert.Nil(t, err)
	rr := httptest.NewRecorder()
	app.ServeHTTP(rr, req)
	then(rr)
}

func TestApp(t *testing.T) {
	t.Run("get query string", func(t *testing.T) {
		type Query struct {
			ID string `query:"id"`
		}
		app := mango.New()
		app.GET("/", func(c *mango.Context) mango.Response {
			query := validate.Query[Query](c)
			return json.OK(map[string]any{"id": query.ID})
		})
		run(t, app, http.MethodGet, "/?id=foo", nil, func(rr *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusOK, rr.Result().StatusCode)
			var resp map[string]any
			assert.Nil(t, stdjson.NewDecoder(rr.Body).Decode(&resp))
			assert.Equal(t, "foo", resp["id"])
		})
	})
	t.Run("get query int", func(t *testing.T) {
		type Query struct {
			ID int `query:"id"`
		}
		app := mango.New()
		app.GET("/", func(c *mango.Context) mango.Response {
			query := validate.Query[Query](c)
			return json.OK(map[string]any{"id": query.ID})
		})
		run(t, app, http.MethodGet, "/?id=5", nil, func(rr *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusOK, rr.Result().StatusCode)
			var resp map[string]int
			assert.Nil(t, stdjson.NewDecoder(rr.Body).Decode(&resp))
			assert.Equal(t, 5, resp["id"])
		})
	})
	t.Run("get query 400", func(t *testing.T) {
		type Query struct {
			ID int `query:"id"`
		}
		app := mango.New()
		app.GET("/", func(c *mango.Context) mango.Response {
			query := validate.Query[Query](c)
			return json.OK(map[string]any{"id": query.ID})
		})
		run(t, app, http.MethodGet, "/?id=foo", nil, func(rr *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusBadRequest, rr.Result().StatusCode)
		})
	})
	t.Run("get query 500 (unexported field)", func(t *testing.T) {
		type Query struct {
			id string `query:"id"`
		}
		app := mango.New()
		app.GET("/", func(c *mango.Context) mango.Response {
			query := validate.Query[Query](c)
			return json.OK(map[string]any{"id": query.id})
		})
		run(t, app, http.MethodGet, "/?id=foo", nil, func(rr *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusInternalServerError, rr.Result().StatusCode)
		})
	})
	t.Run("get path string", func(t *testing.T) {
		type Path struct {
			ID string `path:"id"`
		}
		app := mango.New()
		app.GET("/foo/{id}", func(c *mango.Context) mango.Response {
			path := validate.Path[Path](c)
			return json.OK(map[string]any{"id": path.ID})
		})
		run(t, app, http.MethodGet, "/foo/bar", nil, func(rr *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusOK, rr.Result().StatusCode)
			var resp map[string]any
			assert.Nil(t, stdjson.NewDecoder(rr.Body).Decode(&resp))
			assert.Equal(t, "bar", resp["id"])
		})
	})
	t.Run("get path int", func(t *testing.T) {
		type Path struct {
			ID int `path:"id"`
		}
		app := mango.New()
		app.GET("/foo/{id}", func(c *mango.Context) mango.Response {
			path := validate.Path[Path](c)
			return json.OK(map[string]any{"id": path.ID})
		})
		run(t, app, http.MethodGet, "/foo/5", nil, func(rr *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusOK, rr.Result().StatusCode)
			var resp map[string]int
			assert.Nil(t, stdjson.NewDecoder(rr.Body).Decode(&resp))
			assert.Equal(t, 5, resp["id"])
		})
	})
	t.Run("post body", func(t *testing.T) {
		type Body struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		}
		app := mango.New()
		app.POST("/", func(c *mango.Context) mango.Response {
			return json.OK(validate.Body[Body](c))
		})
		run(t, app, http.MethodPost, "/", &Body{ID: 5, Name: "foo"}, func(rr *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusOK, rr.Result().StatusCode)
			var resp Body
			assert.Nil(t, stdjson.NewDecoder(rr.Body).Decode(&resp))
			assert.Equal(t, 5, resp.ID)
			assert.Equal(t, "foo", resp.Name)
		})
	})
	t.Run("post body 400 invalid", func(t *testing.T) {
		type Body struct {
			ID   int    `json:"id"`
			Name string `json:"name" valid:"req|min:3|max:32"`
		}
		app := mango.New()
		app.POST("/", func(c *mango.Context) mango.Response {
			return json.OK(validate.Body[Body](c))
		})
		run(t, app, http.MethodPost, "/", map[string]any{"name": "Jo"}, func(rr *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusUnprocessableEntity, rr.Result().StatusCode)
		})
	})
	t.Run("post body 400 wrong type", func(t *testing.T) {
		type Body struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		}
		app := mango.New()
		app.POST("/", func(c *mango.Context) mango.Response {
			path := validate.Body[Body](c)
			return json.OK(path)
		})
		run(t, app, http.MethodPost, "/", map[string]any{"Name": 6}, func(rr *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusBadRequest, rr.Result().StatusCode)
		})
	})
	t.Run("custom error handler", func(t *testing.T) {
		type Body struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		}
		app := mango.New(mango.WithErrorHandler(func(c *mango.Context, err error) mango.Response {
			switch {
			case c.Status() >= 500:
				return json.Response(c.Status(), map[string]string{"message": "internal server error"})
			default:
				return json.Response(c.Status(), map[string]string{"message": err.Error()})
			}
		}))
		app.POST("/", func(c *mango.Context) mango.Response {
			body := validate.Body[Body](c)
			return json.OK(body)
		})
		run(t, app, http.MethodPost, "/", map[string]any{"name": 6}, func(rr *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusBadRequest, rr.Result().StatusCode)
			var resp map[string]string
			assert.Nil(t, stdjson.NewDecoder(rr.Body).Decode(&resp))
			assert.Contains(t, resp["message"], "cannot unmarshal")
		})
	})
	t.Run("panic", func(t *testing.T) {
		app := mango.New()
		app.GET("/", func(c *mango.Context) mango.Response {
			v := make([]int, 0)
			return json.OK(v[1])
		})
		run(t, app, http.MethodGet, "/", nil, func(rr *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusInternalServerError, rr.Result().StatusCode)
		})
	})
	t.Run("use middleware", func(t *testing.T) {
		app := mango.New()
		app.Use(func(c *mango.Context) {
			c.Set("foo", "bar")
			c.Next()
		})
		app.GET("/", func(c *mango.Context) mango.Response {
			foo := c.Get("foo")
			return json.OK(map[string]any{"foo": foo})
		})
		run(t, app, http.MethodGet, "/", nil, func(rr *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusOK, rr.Result().StatusCode)
			var resp map[string]any
			assert.Nil(t, stdjson.NewDecoder(rr.Body).Decode(&resp))
			assert.Equal(t, "bar", resp["foo"])
		})
	})
	t.Run("group route used", func(t *testing.T) {
		app := mango.New()
		group := app.Group(mango.WithMiddleware(func(c *mango.Context) {
			c.Set("foo", "bar")
			c.Next()
		}))
		group.GET("/", func(c *mango.Context) mango.Response {
			foo := c.Get("foo")
			return json.OK(map[string]any{"foo": foo})
		})
		run(t, app, http.MethodGet, "/", nil, func(rr *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusOK, rr.Result().StatusCode)
			var resp map[string]any
			assert.Nil(t, stdjson.NewDecoder(rr.Body).Decode(&resp))
			assert.Equal(t, "bar", resp["foo"])
		})
	})
	t.Run("group route unused", func(t *testing.T) {
		app := mango.New()
		group := app.Group(mango.WithMiddleware(func(c *mango.Context) {
			c.Set("foo", "bar")
			c.Next()
		}))
		group.GET("/foo", func(c *mango.Context) mango.Response {
			foo := c.Get("foo")
			return json.OK(map[string]any{"foo": foo})
		})
		app.GET("/", func(c *mango.Context) mango.Response {
			foo := c.Get("foo")
			return json.OK(map[string]any{"foo": foo})
		})
		run(t, app, http.MethodGet, "/", nil, func(rr *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusOK, rr.Result().StatusCode)
			var resp map[string]any
			assert.Nil(t, stdjson.NewDecoder(rr.Body).Decode(&resp))
			assert.Nil(t, resp["foo"])
		})
	})
}
