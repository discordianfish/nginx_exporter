package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	nginxStatus = `Active connections: 91 
server accepts handled requests
 145249 145249 151557 
Reading: 0 Writing: 24 Waiting: 66 
`
	metricCount = 7
)

func TestNginxStatus(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(nginxStatus))
	})
	server := httptest.NewServer(handler)

	e := NewExporter(server.URL)
	ch := make(chan prometheus.Metric)

	go func() {
		defer close(ch)
		e.Collect(ch)
	}()

	for i := 1; i <= metricCount; i++ {
		m := <-ch
		if m == nil {
			t.Error("expected metric but got nil")
		}
	}
	if <-ch != nil {
		t.Error("expected closed channel")
	}
}
