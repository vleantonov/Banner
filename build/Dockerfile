FROM golang:alpine AS builder

LABEL stage=gobuilder

ENV CGO_ENABLED=0
ENV GOOS=linux

RUN apk update --no-cache && apk add --no-cache tzdata

WORKDIR /build
ADD ./go.mod .
ADD ./go.sum .

RUN go mod download

COPY . .

RUN go build -ldflags="-s -w" -o /app/banner ./cmd/banner/main.go
RUN go build -ldflags="-s -w" -o /app/migrator ./cmd/migrator/main.go
RUN go build -ldflags="-s -w" -o /app/worker ./cmd/worker/main.go

FROM alpine

RUN apk update --no-cache && apk add --no-cache ca-certificates

COPY --from=builder /usr/share/zoneinfo/Europe/Moscow /usr/share/zoneinfo/Europe/Moscow

ENV TZ=Europe/Moscow

ENV APP_PORT=8888
ENV APP_HOST=0.0.0.0

COPY --from=builder /app /app

WORKDIR /app

COPY --from=builder /build/migrations ./migrations

CMD ./migrator && ./banner