FROM golang:1.25 AS builder
WORKDIR /app

# Install swag CLI for generating Swagger documentation
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Copy dependency files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Generate Swagger documentation
RUN swag init -g cmd/api/main.go -o docs

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o order-packing-api ./cmd/api

FROM scratch
WORKDIR /app
COPY --from=builder /app/order-packing-api ./order-packing-api
COPY --from=builder /app/static ./static
COPY --from=builder /app/docs ./docs
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
EXPOSE 8080
ENTRYPOINT ["./order-packing-api"]
