FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o actors-service ./cmd

FROM alpine:3.21
WORKDIR /app

COPY --from=builder /app/actors-service .

EXPOSE 8003
CMD ["./actors-service"]