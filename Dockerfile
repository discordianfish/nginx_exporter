FROM	golang:1.9
LABEL	maintainer="@discordianfish"

# Install dep tool
ARG	DEP_VERSION=v0.3.2
ARG	DEP_SHA256=322152b8b50b26e5e3a7f6ebaeb75d9c11a747e64bbfd0d8bb1f4d89a031c2b5
RUN	wget -q https://github.com/golang/dep/releases/download/${DEP_VERSION}/dep-linux-amd64 -O /usr/local/bin/dep \
&&	echo "${DEP_SHA256}  /usr/local/bin/dep" | sha256sum -c - \
&&	chmod 755 /usr/local/bin/dep

ENV	GOOS=linux CGO_ENABLED=0
WORKDIR	/go/src/github.com/discordianfish/nginx_exporter
COPY	Gopkg.* *.go ./
RUN	dep ensure --vendor-only \
&&	go install

FROM	quay.io/prometheus/busybox:glibc
EXPOSE	9113
COPY	--from=0 /go/bin/nginx_exporter /usr/local/bin/
USER	nobody
ENTRYPOINT [ "/usr/local/bin/nginx_exporter" ]
