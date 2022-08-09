module github.com/discordianfish/nginx_exporter

go 1.16

require (
	github.com/prometheus/client_golang v1.0.0
	github.com/prometheus/common v0.6.0
)

replace (
	github.com/golang/protobuf v1.4.0-rc.4.0.20200313231945-b860323f09d0 => github.com/golang/protobuf v1.4.0
	github.com/mwitkow/go-conntrack v0.0.0-20161129095857-cc309e4a2223 => github.com/mwitkow/go-conntrack v0.0.0-20190716064945-2f068394615f
	google.golang.org/protobuf v1.23.1-0.20200526195155-81db48ad09cc => google.golang.org/protobuf v1.24.0
)
