FROM golang:1.9
LABEL maintainer="@discordianfish"
WORKDIR /go/src/github.com/discordianfish/nginx_exporter
ENV GOOS=linux CGO_ENABLED=0
COPY . .
RUN  go get -d && go build


FROM quay.io/prometheus/busybox:glibc
EXPOSE 9113
COPY --from=0 /go/src/github.com/discordianfish/nginx_exporter/nginx_exporter /bin/
USER nobody
ENTRYPOINT [ "/bin/nginx_exporter" ]
