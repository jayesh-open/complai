FROM golang:1.22-alpine AS builder

WORKDIR /build

COPY packages/shared-kernel-go/ ./packages/shared-kernel-go/

COPY services/go/identity-service/go.mod services/go/identity-service/go.sum ./services/go/identity-service/
WORKDIR /build/services/go/identity-service
RUN go mod download

COPY services/go/identity-service/ ./

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/server ./cmd/server

FROM alpine:3.20

RUN apk --no-cache add ca-certificates tzdata
RUN adduser -D -u 1000 appuser

COPY --from=builder /app/server /app/server
COPY services/go/identity-service/migrations /app/migrations

USER appuser
EXPOSE 8081

ENTRYPOINT ["/app/server"]
