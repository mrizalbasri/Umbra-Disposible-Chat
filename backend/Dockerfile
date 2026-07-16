# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy dependency files first for caching
COPY . ./
RUN ls -R

# Production stage
FROM alpine:latest  

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from build stage
COPY --from=builder /app/umbra-server .

# Heroku will dynamically set the PORT environment variable
ENV PORT=8080

CMD ["./umbra-server"]
