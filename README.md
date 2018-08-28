# Nginx Exporter for Prometheus

This Prometheus Exporter retrieves nginx stats and exports them via HTTP for Prometheus
consumption.

To run it:

```bash
./nginx_exporter [flags]
```

Help on flags:
```bash
./nginx_exporter --help
```

## Using Docker

```
docker pull fish/nginx-exporter

docker run -d -p 9113:9113 fish/nginx-exporter \
    -nginx.scrape_uri=http://172.17.42.1/nginx_status
```
In production you should use a tagged release:
https://hub.docker.com/r/fish/nginx-exporter/tags/

## Alternatives
While nginx natively only provides the small set of metrics this exporter
provides, [nginx-module-vts](https://github.com/vozlt/nginx-module-vts)
adds extensive metrics that can be consumed by:

- Standalone Prometheus Exporter: https://github.com/hnlq715/nginx-vts-exporter
- Kubernetes NGINX Ingres controller: https://github.com/kubernetes/ingress-nginx/blob/master/docs/examples/customization/custom-vts-metrics-prometheus/README.md
