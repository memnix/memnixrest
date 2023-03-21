FROM golang:1.20-alpine as builder

LABEL stage=gobuilder

ENV CGO_ENABLED 0
ENV GOOS linux

RUN apk update --no-cache && apk add --no-cache tzdata && apk add --no-cache upx

WORKDIR /build

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN go get -d -v
RUN go build -ldflags="-s -w" -o /app/memnixrest .
RUN upx /app/memnixrest

FROM alpine:3.17

RUN apk update --no-cache && apk add --no-cache ca-certificates
COPY --from=builder /usr/share/zoneinfo/Europe/Paris /usr/share/zoneinfo/Europe/Paris
ENV TZ Europe/Paris

WORKDIR /app

COPY --from=builder /app/memnixrest /app/memnixrest
COPY --from=builder /build/.env /app/.env
COPY --from=builder /build/favicon.ico /app/favicon.ico


EXPOSE 1815

RUN apk add --no-cache dumb-init
ENTRYPOINT ["/usr/bin/dumb-init", "--"]
CMD ["/app/memnixrest"]
