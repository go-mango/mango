package mango_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-mango/mango"
	"github.com/go-mango/mango/json"
)

var plainTextBody = "Hello, World!"

var mangoApp *mango.App

type Post struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

type MangoPath struct {
	Foo string `path:"foo"`
	Bar string `path:"bar"`
	Baz string `path:"baz"`
}

var posts []Post

func init() {
	for i := 0; i < 100; i++ {
		posts = append(posts, Post{
			ID:    1234,
			Title: "Lorem Ipsum",
			Body:  "Veniam ipsum officia consequat minim veniam cillum incididunt laborum aliqua ad do magna sed aliquip fugiat. Cillum et aliqua commodo, velit minim anim, pariatur, magna culpa officia dolor quis consectetur. Proident commodo laboris eu eu quis esse ea exercitation irure pariatur duis nulla deserunt dolor sed. Nulla qui laboris ut ea non consectetur amet culpa, pariatur commodo magna deserunt nostrud in.",
		})
	}

	mangoApp = mango.New()
	// mangoApp.GET("/plaintext", func(c *mango.Context) mango.Responder {
	// 	return c.Text(http.StatusOK, plainTextBody)
	// })
	mangoApp.GET("/json", func(c *mango.Context) mango.Response {
		return json.OK(&posts)
	})
	mangoApp.GET("/params/{foo}/{bar}/{baz}", func(c *mango.Context) mango.Response {
		// return json.OK(mango.ParsePath[MangoPath](c))
		path := mango.ParsePath[MangoPath](c)
		return json.OK(map[string]any{
			"foo": path.Foo,
			"bar": path.Bar,
			"baz": path.Baz,
		})
		// path := mango.ParsePath[MangoPath](c)
		// return json.OK(map[string]any{
		// 	"foo": c.Request().PathValue("foo"),
		// 	"bar": c.Request().PathValue("bar"),
		// 	"baz": c.Request().PathValue("baz"),
		// })
	})
}

type Fatalfer interface {
	Fatalf(string, ...interface{})
}

func equals(f Fatalfer, a interface{}, b interface{}) {
	if a != b {
		f.Fatalf("expected %v to equal %v", a, b)
	}
}

// BenchmarkMangoJSON-16    	   31105	     38594 ns/op	   50748 B/op	      11 allocs/op
func BenchmarkMangoJSON(b *testing.B) {
	req, err := http.NewRequest("GET", "/json", nil)
	equals(b, err, nil)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		mangoApp.ServeHTTP(rr, req)
		equals(b, rr.Code, http.StatusOK)
	}
}

func BenchmarkMangoPathParams(b *testing.B) {
	req, err := http.NewRequest("GET", "/params/a/b/c", nil)
	equals(b, err, nil)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		mangoApp.ServeHTTP(rr, req)
		equals(b, rr.Code, http.StatusOK)
	}
}
