FROM golang:1.5

ENTRYPOINT /go/bin/go-reverse-proxy

EXPOSE 80

RUN go get gopkg.in/yaml.v2

COPY . /go/src/github.com/wouterd/go-reverse-proxy

RUN go install github.com/wouterd/go-reverse-proxy
