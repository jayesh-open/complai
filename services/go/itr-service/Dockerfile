FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY packages/shared-kernel-go/ packages/shared-kernel-go/
COPY services/go/itr-service/ services/go/itr-service/
WORKDIR /app/services/go/itr-service
RUN go mod download
RUN CGO_ENABLED=0 go build -o /itr-service ./cmd/server

FROM alpine:3.19
RUN apk --no-cache add ca-certificates
COPY --from=builder /itr-service /itr-service
COPY --from=builder /app/services/go/itr-service/migrations /migrations
EXPOSE 8100
ENTRYPOINT ["/itr-service"]
