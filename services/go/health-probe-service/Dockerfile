FROM golang:1.22-alpine AS builder

WORKDIR /build

COPY packages/shared-kernel-go/ /build/packages/shared-kernel-go/

COPY services/go/health-probe-service/go.mod services/go/health-probe-service/go.sum ./
RUN go mod download

COPY services/go/health-probe-service/ ./

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/server ./cmd/server

FROM alpine:3.20

RUN apk --no-cache add ca-certificates tzdata
RUN adduser -D -u 1000 appuser

COPY --from=builder /app/server /app/server

USER appuser
EXPOSE 8080

ENTRYPOINT ["/app/server"]
