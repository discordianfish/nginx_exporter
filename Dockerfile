FROM       ubuntu
MAINTAINER Johannes 'fish' Ziemke <github@freigeist.org> @discordianfish

RUN        apt-get update && apt-get install -yq curl git mercurial gcc
RUN        curl -s https://go.googlecode.com/files/go1.2.linux-amd64.tar.gz | tar -C /usr/local -xzf -
ENV        PATH    /usr/local/go/bin:$PATH
ENV        GOPATH  /go

ADD        . /usr/src/nginx_exporter
RUN        cd /usr/src/nginx_exporter && \
           go get -d && go build && cp nginx_exporter /

ENTRYPOINT [ "/nginx_exporter" ]
EXPOSE     8080
