# Mango
Simple web framework.


## Basic Usage
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