# Базовый образ
FROM golang:1.24-alpine AS builder


RUN apk add --no-cache git


WORKDIR /app


COPY go.mod go.sum ./
RUN go mod download


COPY . .


RUN CGO_ENABLED=0 GOOS=linux go build -o api-gateway ./cmd


FROM alpine:3.21
WORKDIR /app


COPY --from=builder /app/api-gateway .

EXPOSE 8080

CMD ["./api-gateway"]