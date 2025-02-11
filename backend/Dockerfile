FROM golang:alpine AS builder

WORKDIR /build
ADD go.mod .
COPY .. .
RUN go build -o app ./cmd/markup/markup.go

FROM alpine

WORKDIR /app
COPY --from=builder /build/app /app/app
COPY ./config ./config
RUN ls -la
RUN "pwd"
CMD ["./app", "--config=./config/local.yml"]