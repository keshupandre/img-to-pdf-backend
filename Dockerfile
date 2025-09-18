# Use official Go image
FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY . .

# Build the Go binary
RUN go build -o app ./cmd

# Final lightweight image
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/app .

EXPOSE 8080
CMD ["./app"]