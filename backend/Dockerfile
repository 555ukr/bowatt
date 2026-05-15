# Build stage
FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/bowatt cmd/main.go

# Run stage
FROM alpine:3.20

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /app/bin/bowatt .

RUN mkdir -p /app/uploads

EXPOSE 8000

ENTRYPOINT ["./bowatt"]
