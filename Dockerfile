# Use a lightweight base image with Go
FROM golang:1.20-alpine AS builder

# Set up environment variables
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Install required tools
RUN apk add --no-cache git

# Set the working directory
WORKDIR /app

# Copy the Go module files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the binary
RUN go build -o openapi-axe .

# Use a lightweight runtime image for the final container
FROM alpine:3.18

# Set up a non-root user for better security
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Set the working directory
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/openapi-axe /usr/local/bin/openapi-axe

# Set permissions
RUN chown -R appuser:appgroup /app

# Switch to the non-root user
USER appuser

# Set the default entrypoint
ENTRYPOINT ["openapi-axe"]
