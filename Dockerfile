FROM alpine:latest

ENV GOPATH /go
ENV APPPATH $GOPATH/src/github.com/discordianfish/nginx_exporter

COPY . $APPPATH

RUN apk add --update -t build-deps go git mercurial \
    && cd $APPPATH && go get -d && go build -o /nginx_exporter \
    && apk del --purge build-deps git mercurial && rm -rf $GOPATH

EXPOSE 9113

ENTRYPOINT ["/nginx_exporter"]
