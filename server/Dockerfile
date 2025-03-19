FROM golang:1.24-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY internal internal
COPY *.go ./
RUN go build -v -o app

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/app .
EXPOSE 80
CMD ["./app"]