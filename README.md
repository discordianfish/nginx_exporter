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

## Using Docker

```
docker pull fish/nginx-exporter

docker run -d -p 9113:9113 fish/nginx-exporter \
    -nginx.scrape_uri=http://172.17.42.1/nginx_status
```
