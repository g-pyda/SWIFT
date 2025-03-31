FROM golang:1.24.1 AS builder
ENV GOOS=linux GOARCH=amd64 CGO_ENABLED=0

WORKDIR /build
COPY ./go.mod ./go.sum .
RUN go mod download

COPY ./src ./data .
RUN go build -o ./swift ./src/main.go

FROM alpine:latest AS production

WORKDIR /app
COPY --from=builder /build/swift ./swift
RUN chmod +x /app/swift

COPY --from=builder /build/data/Interns_2025_SWIFT_CODES.xlsx ./data/

CMD ["./swift"]