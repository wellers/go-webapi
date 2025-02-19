FROM golang:1.24-alpine
EXPOSE 80
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./
RUN go build -v -o /usr/local/bin/app ./...

CMD ["app"]