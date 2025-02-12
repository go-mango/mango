module github.com/go-mango/mango

go 1.23.3

require (
	github.com/go-mango/json v0.0.0
	github.com/go-mango/validate v0.0.0
	github.com/stretchr/testify v1.10.0
	github.com/twharmon/govalid v1.5.1 // indirect
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/go-mango/validate v0.0.0 => ../validate

replace github.com/go-mango/json v0.0.0 => ../json

replace github.com/go-mango/dynamodb v0.0.0 => ../dynamodb
