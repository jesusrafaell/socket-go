FROM golang:alpine AS builder

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o bin cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/bin .

EXPOSE 8080

CMD ["./bin"]
