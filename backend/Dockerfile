ARG BUILD_PATH=./cmd/markup/markup.go

FROM golang:alpine AS builder

WORKDIR /build
ADD go.mod .
COPY .. .

ARG BUILD_PATH
RUN go build -o app ${BUILD_PATH}

FROM alpine

WORKDIR /app
COPY --from=builder /build/app /app/app
COPY ./config ./config

CMD ["./app", "--config=./config/dev.yml"]