FROM alpine:3.2
MAINTAINER The Prometheus Authors <prometheus-developers@googlegroups.com>

ENV GOPATH /go
COPY . /go/src/github.com/prometheus/nginx_exporter

RUN apk add --update -t build-deps go git mercurial make \
    && apk add -u musl && rm -rf /var/cache/apk/* \
    && cd /go/src/github.com/prometheus/nginx_exporter \
    && make && cp nginx_exporter /bin/nginx_exporter \
    && rm -rf /go && apk del --purge build-deps

EXPOSE     9113
ENTRYPOINT [ "/bin/nginx_exporter" ]
