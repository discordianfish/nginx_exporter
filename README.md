# Nginx Exporter for Prometheus

This is a simple server that periodically scrapes nginx stats and exports them via HTTP for Prometheus
consumption.

To run it:

```bash
./nginx_exporter [flags]
```

Help on flags:
```bash
./nginx_exporter --help
```

## Getting Started
  * All of the core developers are accessible via the [Prometheus Developers Mailinglist](https://groups.google.com/forum/?fromgroups#!forum/prometheus-developers).

## Building

From scratch (on Ubuntu/Debian)
```bash
apt-cache search golang
```
Latest version as of this writing was 1.9
```bash
sudo apt-get install golang-1.9
cd ~
mkdir -p go/bin
```
You might want to add the following 3 lines to your .bashrc to make it permanent
```bash
export PATH=$PATH:/usr/lib/go-1.9/bin
export GOPATH=$HOME/go
export GOBIN=$HOME/go/bin
```
Clone, get, build: done
```bash
git clone https://github.com/discordianfish/nginx_exporter.git
cd nginx_exporter/
go get
go build
```

## Using Docker

```
docker pull fish/nginx-exporter

docker run -d -p 9113:9113 fish/nginx-exporter \
    -nginx.scrape_uri=http://172.17.42.1/nginx_status
```
