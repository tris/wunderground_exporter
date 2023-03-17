# Builder stage
FROM golang:1.17-alpine as builder

MAINTAINER Tristan Horn <tristan+docker@ethereal.net>

WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the entire source code
COPY . .

# Compile the Go code
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o wunderground_exporter .

# Final stage
FROM scratch

# Copy the compiled binary from the builder stage
COPY --from=builder /app/wunderground_exporter /wunderground_exporter

# Add the CA certificates to trust HTTPS connections
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

EXPOSE 9122

CMD ["/wunderground_exporter"]
