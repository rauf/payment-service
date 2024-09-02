FROM golang:1.23-alpine AS builder
RUN apk update && apk add --no-cache make

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN make build

FROM alpine:latest

WORKDIR /app
ENV ENV=prod

COPY --from=builder /app/tmp/api ./bin/api
EXPOSE 3000
ENTRYPOINT ["./bin/api"]
