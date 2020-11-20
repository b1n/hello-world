FROM golang:1.15-alpine as builder

ENV GO111MODULE on

WORKDIR /go/src/hello-world

COPY . /go/src/hello-world

RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh && \
    go get -d -v ./... && go install -v

FROM alpine:3.10

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=builder /go/bin/hello-world /app/app

CMD ["/app/app"]

EXPOSE 9111
