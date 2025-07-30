FROM golang:alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . /app
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/memproxy ./cmd/memproxy

# Deploy the application binary into a lean image
FROM alpine AS build-release-stage
WORKDIR /app
COPY --from=builder  /app/memproxy ./memproxy
EXPOSE 8080
CMD [ "--listen=8080", "--backend=memcached", "--origin=http://localhost:8081" ]
ENTRYPOINT ["/app/memproxy"]