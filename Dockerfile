# syntax=docker/dockerfile:1
FROM golang:1.23-alpine AS build

WORKDIR /app

# Install wget for downloading Tailwind CLI
RUN apk add --no-cache wget

# Download Tailwind CLI
RUN wget -O /usr/local/bin/tailwindcss https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-linux-x64 && \
    chmod +x /usr/local/bin/tailwindcss

# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build Tailwind CSS
RUN /usr/local/bin/tailwindcss -i ./assets/web/css/tailwind.css -o ./assets/web/css/app.css --minify

# Build with caching
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

# Final stage
FROM alpine:latest

WORKDIR /app

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Copy binary from build stage
COPY --from=build /app/server .

EXPOSE 8080

CMD ["./server"]
