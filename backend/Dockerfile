FROM golang:alpine AS builder

WORKDIR /build
ADD go.mod .
COPY .. .
RUN go build -o app ./cmd/markup/markup.go

FROM alpine

WORKDIR /app
COPY --from=builder /build/app /app/app
COPY ./config ./config

CMD ["./app", "--config=./config/config.yml"]