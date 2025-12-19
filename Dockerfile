FROM golang:1.24.5 AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -o ./hello-app ./cmd/main.go

FROM scratch AS final
COPY --from=builder /app/hello-app /hello-app
COPY --from=builder /app/.env /.env
COPY --from=builder /app/internal/migrations /internal/migrations

EXPOSE "8080"
ENTRYPOINT ["/hello-app"]