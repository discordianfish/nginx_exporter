FROM golang:1.9
WORKDIR /go/src/github.com/discordianfish/nginx_exporter
COPY . .
RUN set -xe; \
  go get -d; \
  CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo;


FROM quay.io/prometheus/busybox
EXPOSE 9113
COPY --from=0 /go/src/github.com/discordianfish/nginx_exporter/nginx_exporter /bin/
USER nobody
ENTRYPOINT [ "/bin/nginx_exporter" ]
