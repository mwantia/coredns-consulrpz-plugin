# Build stage
FROM golang:1.22 AS builder

WORKDIR /app

# Copy the local plugin files
COPY . /app/coredns-rpz-plugin

# Clone CoreDNS repository
RUN git clone https://github.com/coredns/coredns.git
WORKDIR /app/coredns

# Replace plugin.cfg
COPY ./build/plugin.cfg /app/coredns/plugin.cfg

# Update go.mod to use the local plugin
RUN go mod edit -replace github.com/mwantia/coredns-rpz-plugin=/app/coredns-rpz-plugin

# Update dependencies and build
RUN go get github.com/mwantia/coredns-env-plugin
RUN go get github.com/mwantia/coredns-rpz-plugin
RUN go generate
RUN go mod tidy
RUN make

# Final stage
FROM debian:bullseye-slim

WORKDIR /app

# Copy the built CoreDNS binary from the builder stage
COPY --from=builder /app/coredns/coredns /app/coredns

# Expose DNS ports
EXPOSE 53/udp
EXPOSE 53/tcp

# Run CoreDNS
CMD ["/app/coredns"]