# Test Stage (Separate)
FROM golang:1.24-alpine AS tester
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY internal internal
COPY *.go ./

CMD ["go", "test", "-v"]