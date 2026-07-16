# Build stage
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Copy dependency files first for caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go application binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o umbra-server cmd/server/main.go

# Production stage
FROM alpine:latest  

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from build stage
COPY --from=builder /app/umbra-server .

# Heroku will dynamically set the PORT environment variable
ENV PORT=8080

CMD ["./umbra-server"]
