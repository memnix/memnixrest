FROM golang:1.21-alpine AS builder

LABEL stage=gobuilder

ENV CGO_ENABLED 0
ENV GOOS linux

RUN apk update --no-cache && apk add --no-cache tzdata && apk add --no-cache upx

WORKDIR /build

COPY go.mod .
COPY go.sum .
COPY default.pgo .
RUN go mod download

COPY . .

RUN go get -d -v \
    && go build -pgo=auto -ldflags="-s -w" -o /app/memnixrest .\
    && upx /app/memnixrest

FROM alpine:3.18

RUN addgroup -S memnix && adduser -S memnix -G memnix \
    && apk update --no-cache && apk add --no-cache ca-certificates

COPY --from=builder /usr/share/zoneinfo/Europe/Paris /usr/share/zoneinfo/Europe/Paris
ENV TZ Europe/Paris

ENV GOMEMLIMIT 4000MiB

WORKDIR /app

RUN mkdir -p /app/config/keys

COPY --from=builder /app/memnixrest /app/memnixrest
COPY --from=builder /build/.env* /app/.
COPY --from=builder /build/favicon.ico /app/favicon.ico
COPY --from=builder /build/config/keys /app/config/keys

# Change ownership of the app directory to the non-root user
RUN chown -R memnix:memnix /app

EXPOSE 1815

RUN apk add --no-cache dumb-init
USER memnix
ENTRYPOINT ["/usr/bin/dumb-init", "--"]
CMD ["/app/memnixrest"]
