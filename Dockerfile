FROM golang:1.5

COPY . /go/src/github.com/wouterd/go-reverse-proxy

RUN go install github.com/wouterd/go-reverse-proxy

ENTRYPOINT /go/bin/go-reverse-proxy

EXPOSE 80
