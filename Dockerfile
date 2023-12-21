# Build
FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY src/ .

RUN CGO_ENABLED=0 GOOS=linux go build -o main


# Run
FROM alpine
WORKDIR /opt
COPY .env .
COPY --from=builder /app/main .

CMD [ "./main" ]