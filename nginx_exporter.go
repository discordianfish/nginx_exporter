package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

const (
	namespace = "nginx" // For Prometheus metrics.
)

var (
	listeningAddress = flag.String("telemetry.address", ":9113", "Address on which to expose metrics.")
	metricsEndpoint  = flag.String("telemetry.endpoint", "/metrics", "Path under which to expose metrics.")
	nginxScrapeURI   = flag.String("nginx.scrape_uri", "http://localhost/nginx_status", "URI to nginx stub status page")
	insecure         = flag.Bool("insecure", true, "Ignore server certificate if using https")
)

// Exporter collects nginx stats from the given URI and exports them using
// the prometheus metrics package.
type Exporter struct {
	URI    string
	mutex  sync.RWMutex
	client *http.Client

	scrapeFailures       prometheus.Counter
	processedConnections *prometheus.Desc
	currentConnections   *prometheus.GaugeVec
	nginxUp              prometheus.Gauge
}

// NewExporter returns an initialized Exporter.
func NewExporter(uri string) *Exporter {
	return &Exporter{
		URI: uri,
		scrapeFailures: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "exporter_scrape_failures_total",
			Help:      "Number of errors while scraping nginx.",
		}),
		processedConnections: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "connections_processed_total"),
			"Number of connections processed by nginx",
			[]string{"stage"},
			nil,
		),
		currentConnections: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "connections_current",
			Help:      "Number of connections currently processed by nginx",
		},
			[]string{"state"},
		),
		nginxUp: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "up",
			Help:      "Whether the nginx is up.",
		}),
		client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: *insecure},
			},
		},
	}
}

// Describe describes all the metrics ever exported by the nginx exporter. It
// implements prometheus.Collector.
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- e.processedConnections
	e.currentConnections.Describe(ch)
	e.nginxUp.Describe(ch)
	e.scrapeFailures.Describe(ch)
}

func (e *Exporter) collect(ch chan<- prometheus.Metric) error {
	resp, err := e.client.Get(e.URI)
	if err != nil {
		e.nginxUp.Set(0)
		return fmt.Errorf("Error scraping nginx: %v", err)
	}
	e.nginxUp.Set(1)

	data, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		if err != nil {
			data = []byte(err.Error())
		}
		return fmt.Errorf("Status %s (%d): %s", resp.Status, resp.StatusCode, data)
	}

	// Parsing results
	lines := strings.Split(string(data), "\n")
	if len(lines) != 5 {
		return fmt.Errorf("Unexpected number of lines in status: %v", lines)
	}

	// active connections
	parts := strings.Split(lines[0], ":")
	if len(parts) != 2 {
		return fmt.Errorf("Unexpected first line: %s", lines[0])
	}
	v, err := strconv.Atoi(strings.TrimSpace(parts[1]))
	if err != nil {
		return err
	}
	e.currentConnections.WithLabelValues("active").Set(float64(v))

	// processed connections
	parts = strings.Fields(lines[2])
	if len(parts) != 3 {
		return fmt.Errorf("Unexpected third line: %s", lines[2])
	}
	v, err = strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return err
	}
	ch <- prometheus.MustNewConstMetric(e.processedConnections, prometheus.CounterValue, float64(v), "accepted")
	v, err = strconv.Atoi(strings.TrimSpace(parts[1]))
	if err != nil {
		return err
	}
	ch <- prometheus.MustNewConstMetric(e.processedConnections, prometheus.CounterValue, float64(v), "handled")
	v, err = strconv.Atoi(strings.TrimSpace(parts[2]))
	if err != nil {
		return err
	}
	ch <- prometheus.MustNewConstMetric(e.processedConnections, prometheus.CounterValue, float64(v), "any")

	// current connections
	parts = strings.Fields(lines[3])
	if len(parts) != 6 {
		return fmt.Errorf("Unexpected fourth line: %s", lines[3])
	}
	v, err = strconv.Atoi(strings.TrimSpace(parts[1]))
	if err != nil {
		return err
	}
	e.currentConnections.WithLabelValues("reading").Set(float64(v))
	v, err = strconv.Atoi(strings.TrimSpace(parts[3]))
	if err != nil {
		return err
	}

	e.currentConnections.WithLabelValues("writing").Set(float64(v))
	v, err = strconv.Atoi(strings.TrimSpace(parts[5]))
	if err != nil {
		return err
	}
	e.currentConnections.WithLabelValues("waiting").Set(float64(v))
	return nil
}

// Collect fetches the stats from configured nginx location and delivers them
// as Prometheus metrics. It implements prometheus.Collector.
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	e.mutex.Lock() // To protect metrics from concurrent collects.
	defer e.mutex.Unlock()
	if err := e.collect(ch); err != nil {
		log.Errorf("Error scraping nginx: %s", err)
		e.scrapeFailures.Inc()
		e.scrapeFailures.Collect(ch)
	}
	e.currentConnections.Collect(ch)
	e.nginxUp.Collect(ch)
	return
}

func main() {
	flag.Parse()

	exporter := NewExporter(*nginxScrapeURI)
	prometheus.MustRegister(exporter)

	log.Infof("Starting Server: %s", *listeningAddress)
	http.Handle(*metricsEndpoint, prometheus.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			<head><title>Nginx Exporter</title></head>
			<body>
			<h1>Nginx Exporter</h1>
			<p><a href="` + *metricsEndpoint + `">Metrics</a></p>
			</body>
			</html>`))
	})

	log.Fatal(http.ListenAndServe(*listeningAddress, nil))
}
