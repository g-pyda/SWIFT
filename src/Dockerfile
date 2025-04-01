FROM golang:1.24.1 AS builder
ENV GOOS=linux GOARCH=amd64 CGO_ENABLED=0

WORKDIR /build
COPY ./go.mod ./go.sum .
RUN go mod download

COPY ./src ./src
COPY ./data ./data
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init -g ./src/main.go

RUN go build -o ./swift ./src/main.go

FROM alpine:latest AS production

WORKDIR /app
COPY --from=builder /build/swift ./swift
COPY --from=builder /build/data ./data

RUN chmod +x /app/swift

CMD ["./swift"]
