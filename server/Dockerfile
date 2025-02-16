FROM golang:1.24-alpine

WORKDIR /usr/src/app

COPY go.mod ./
RUN go mod download

COPY . .
RUN go build -v -o /usr/local/bin/app ./...

EXPOSE 80
CMD ["app"]