module github.com/complai/complai/services/go/user-role-service

go 1.22.0

require (
	github.com/complai/complai/packages/shared-kernel-go v0.1.0
	github.com/go-chi/chi/v5 v5.2.1
	github.com/google/uuid v1.6.0
	github.com/jackc/pgx/v5 v5.6.0
	github.com/pressly/goose/v3 v3.21.1
	github.com/rs/zerolog v1.33.0
)

require (
	github.com/go-chi/cors v1.2.1 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/puddle/v2 v2.2.1 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mfridman/interpolate v0.0.2 // indirect
	github.com/rs/xid v1.5.0 // indirect
	github.com/sethvargo/go-retry v0.2.4 // indirect
	github.com/stretchr/testify v1.10.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/crypto v0.24.0 // indirect
	golang.org/x/sync v0.7.0 // indirect
	golang.org/x/sys v0.29.0 // indirect
	golang.org/x/text v0.16.0 // indirect
)

replace github.com/complai/complai/packages/shared-kernel-go => ../../../packages/shared-kernel-go
