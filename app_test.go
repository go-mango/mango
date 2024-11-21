package mango_test

import (
	"bytes"
	stdjson "encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-mango/mango"
	"github.com/go-mango/mango/json"
	"github.com/stretchr/testify/assert"
	"github.com/twharmon/govalid"
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
			query := mango.ParseQuery[Query](c)
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
			query := mango.ParseQuery[Query](c)
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
			query := mango.ParseQuery[Query](c)
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
			query := mango.ParseQuery[Query](c)
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
			path := mango.ParsePath[Path](c)
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
			path := mango.ParsePath[Path](c)
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
			return json.OK(mango.ParseBody[Body](c))
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
		app := mango.New(mango.WithValidator(govalid.Validate))
		app.POST("/", func(c *mango.Context) mango.Response {
			return json.OK(mango.ParseBody[Body](c))
		})
		run(t, app, http.MethodPost, "/", map[string]any{"name": "Jo"}, func(rr *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusBadRequest, rr.Result().StatusCode)
		})
	})
	t.Run("post body 400 wrong type", func(t *testing.T) {
		type Body struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		}
		app := mango.New()
		app.POST("/", func(c *mango.Context) mango.Response {
			path := mango.ParseBody[Body](c)
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
			body := mango.ParseBody[Body](c)
			return json.OK(body)
		})
		run(t, app, http.MethodPost, "/", map[string]any{"name": 6}, func(rr *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusBadRequest, rr.Result().StatusCode)
			var resp map[string]string
			assert.Nil(t, stdjson.NewDecoder(rr.Body).Decode(&resp))
			assert.Contains(t, resp["message"], "cannot unmarshal")
		})
	})
}
