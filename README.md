# Mango
Simple web framework.


## Basic Usage
Path and path parameter matching is handled by the go standard library.

```go
package main

import (
	"github.com/go-mango/mango"
	"github.com/go-mango/mango/json"
)

func main() {
	app := mango.New()

	app.GET("/", func(c *mango.Context) mango.Response {
		return json.OK(map[string]any{
			"hello": "world",
		})
	})

	app.Listen()
}
```

## Input validation
Path parameters, query parameters, and request body can be automatically validated when you provide a validating function.

If mango is unable to parse the provided input, the default (yet customizable) behavior is to return a 400 response with a message: `{"message":"expected 'age' to be an integer"}`.

If validation fails, the default (yet customizable) behavior is to return a 422 response with a message: `{"message":"expected 'age' to be 0-100; got 123"}`.

Internally, mango uses `panic` / `recover` to handle these cases. Although controversial and not idiomatic go, this decision greatly reduces the amount of code required to write http handlers. You can leverage mango's panic recovery mechanism by calling `mango.Abort(statusCode, err)`.

### Path parameters
```go
package main

import (
	"github.com/go-mango/mango"
	"github.com/go-mango/mango/json"
	"github.com/twharmon/govalid"
)

type Params struct {
    Name string `path:"name" valid:"req|max:32"`
}

func main() {
    app := mango.New(mango.WithValidator(govalid.Validate))

	app.GET("/hello/{name}", func(c *mango.Context) mango.Response {
		params := mango.ParsePath[Params](c)
		return json.OK(map[string]any{
			"hello": params.Name,
		})
	})

	app.Listen()

    // GET /hello/Gopher => {"hello":"Gopher"}
}
```

### Query parameters
```go
package main

import (
	"github.com/go-mango/mango"
	"github.com/go-mango/mango/json"
	"github.com/twharmon/govalid"
)

type Params struct {
    Name string `query:"name" valid:"req|max:32"`
}

func main() {
    app := mango.New(mango.WithValidator(govalid.Validate))

	app.GET("/hello", func(c *mango.Context) mango.Response {
		params := mango.ParseQuery[Params](c)
		return json.OK(map[string]any{
			"hello": params.Name,
		})
	})

	app.Listen()

    // GET /hello?name=Gopher => {"hello":"Gopher"}
}
```

### Request body
```go
package main

import (
	"github.com/go-mango/mango"
	"github.com/go-mango/mango/json"
	"github.com/twharmon/govalid"
)

type Params struct {
    Name string `json:"name" valid:"req|max:32"`
}

func main() {
    app := mango.New(mango.WithValidator(govalid.Validate))

	app.POST("/hello", func(c *mango.Context) mango.Response {
		params := mango.ParseBody[Params](c)
		return json.OK(map[string]any{
			"hello": params.Name,
		})
	})

	app.Listen()

    // POST /hello {"name":"Gopher"} => {"hello":"Gopher"}
}
```

## Middleware

Middleware functions are used to process requests before they reach the final handler. They can be used for logging, authentication, or other pre-processing tasks. Here is an example of how to use middleware in your application:

```go
package main

import (
    "fmt"
    "time"

    "github.com/go-mango/mango"
    "github.com/go-mango/mango/json"
)

func main() {
    app := mango.New()

    // Middleware to log the duration of each request
    app.Use(func(c *mango.Context) {
        start := time.Now()
        c.Next()
        fmt.Printf("request took %s\n", time.Since(start))
    })

    app.GET("/", func(c *mango.Context) mango.Response {
        return json.OK(map[string]any{
            "hello": "world",
        })
    })

    app.Listen()
}
```

In this example, the middleware logs the duration of each request. The `c.Next()` function is called to pass control to the next handler in the chain.
