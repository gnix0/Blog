FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/api ./cmd/api
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/migrate ./cmd/migrate

FROM alpine:3.21
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /app/bin/api ./api
COPY --from=builder /app/bin/migrate ./migrate
COPY cmd/migrate/migrations ./cmd/migrate/migrations
RUN adduser -D -u 1001 appuser && chown -R appuser:appuser /app
USER appuser
EXPOSE 8080
CMD ["./api"]
