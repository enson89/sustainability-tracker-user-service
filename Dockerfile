# Build stage
FROM golang:1.20 as builder
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o sustainability-tracker-user-service ./cmd

# Run stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/sustainability-tracker-user-service .
EXPOSE 8080
CMD ["./sustainability-tracker-user-service"]