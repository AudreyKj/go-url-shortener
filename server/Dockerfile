FROM golang:1.24.3-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o url-shortener main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/url-shortener .
COPY --from=builder /app/handlers ./handlers
COPY --from=builder /app/middleware ./middleware
COPY --from=builder /app/services ./services
COPY --from=builder /app/storage ./storage
COPY --from=builder /app/utils ./utils
EXPOSE 8080
CMD ["./url-shortener"]
