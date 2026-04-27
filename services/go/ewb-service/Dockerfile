FROM golang:1.22-alpine AS builder

WORKDIR /build

COPY packages/shared-kernel-go/ ./packages/shared-kernel-go/

COPY services/go/ewb-service/go.mod services/go/ewb-service/go.sum ./services/go/ewb-service/
WORKDIR /build/services/go/ewb-service
RUN go mod download

COPY services/go/ewb-service/ ./

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/server ./cmd/server

FROM alpine:3.20

RUN apk --no-cache add ca-certificates tzdata
RUN adduser -D -u 1000 appuser

COPY --from=builder /app/server /app/server

USER appuser
EXPOSE 8097

ENTRYPOINT ["/app/server"]
