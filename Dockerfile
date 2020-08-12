FROM golang:alpine as builder

COPY . $GOPATH/src/github.com/edmundgmail/ytdl/
RUN apk update && apk upgrade && \
    apk add --no-cache git
RUN go get -d github.com/edmundgmail/ytdl
RUN apk del git
WORKDIR $GOPATH/src/github.com/edmundgmail/ytdl/
RUN go build -o /go/bin/ytdl


FROM alpine
COPY --from=builder /go/bin/ytdl /go/bin/ytdl
COPY --from=builder /go/src/github.com/edmundgmail/ytdl/public /public

ENTRYPOINT ["/go/bin/ytdl"]
