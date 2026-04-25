FROM golang:1.22-alpine AS builder

WORKDIR /build

COPY packages/shared-kernel-go/ ./packages/shared-kernel-go/

COPY services/go/rules-engine-service/go.mod services/go/rules-engine-service/go.sum ./services/go/rules-engine-service/
WORKDIR /build/services/go/rules-engine-service
RUN go mod download

COPY services/go/rules-engine-service/ ./

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/server ./cmd/server

FROM alpine:3.20

RUN apk --no-cache add ca-certificates tzdata
RUN adduser -D -u 1000 appuser

COPY --from=builder /app/server /app/server
COPY services/go/rules-engine-service/migrations /app/migrations

USER appuser
EXPOSE 8090

ENTRYPOINT ["/app/server"]
