FROM golang:1.24-alpine AS tester
COPY --from=go-webapi-shared:latest /app/ /shared
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./

CMD ["go", "test", "-v"]