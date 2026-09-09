# Stage 1: Build
FROM golang:1.26-alpine AS builder

# Install SSL certificates (required for HTTPS requests to ImgBB)
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Compile the binary statically (without OS CGo dependencies)
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o imgbb-service .

# Stage 2: Ultra-lightweight Final Image
FROM scratch

# Copy the SSL certificates from the build stage
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy the compiled binary
COPY --from=builder /app/imgbb-service /imgbb-service

# Expose the default port
EXPOSE 8080

# Run the service
ENTRYPOINT ["/imgbb-service"]
