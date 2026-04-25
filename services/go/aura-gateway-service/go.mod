module github.com/complai/complai/services/go/aura-gateway-service

go 1.23.0

require (
	github.com/complai/complai/packages/shared-kernel-go v0.1.0
	github.com/go-chi/chi/v5 v5.2.1
	github.com/google/uuid v1.6.0
	github.com/rs/zerolog v1.34.0
	github.com/shopspring/decimal v1.4.0
	github.com/stretchr/testify v1.10.0
)

require (
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/go-chi/cors v1.2.1 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/rogpeppe/go-internal v1.14.1 // indirect
	github.com/rs/xid v1.6.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/complai/complai/packages/shared-kernel-go => ../../../packages/shared-kernel-go
