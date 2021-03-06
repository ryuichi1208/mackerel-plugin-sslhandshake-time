package ssl

import (
	"context"
	"crypto/tls"
	"flag"
	"log"
	"net"
	"net/http"
	"net/http/httptrace"
	"net/url"
	"strings"
	"time"

	mp "github.com/mackerelio/go-mackerel-plugin"
)

type Plugin struct {
	Prefix string
	URL    *url.URL
}

type Result struct {
	dnsLookupTime    float64
	tcphandshakeTime float64
	sslhandshakeTime float64
}

func (s Plugin) GraphDefinition() map[string]mp.Graphs {
	labelPrefix := strings.Title(s.MetricKeyPrefix())
	return map[string]mp.Graphs{
		"": {
			Label: labelPrefix + "_connection_time",
			Unit:  mp.UnitFloat,
			Metrics: []mp.Metrics{
				{Name: "dnsLookupTime", Label: "dnsLookupTime"},
				{Name: "tcphandshakeTime", Label: "tcphandshakeTime"},
				{Name: "sslhandshakeTime", Label: "sslhandshakeTime"},
			},
		},
	}
}

func (s Plugin) MetricKeyPrefix() string {
	if s.Prefix == "" {
		s.Prefix = "ssl"
	}
	return s.Prefix
}

func parseURL(uri string) *url.URL {
	if !strings.Contains(uri, "://") && !strings.HasPrefix(uri, "//") {
		uri = "//" + uri
	}

	url, err := url.Parse(uri)
	if err != nil {
		log.Fatalf("could not parse url %q: %v", uri, err)
	}

	if url.Scheme == "" {
		url.Scheme = "https"
	}
	return url
}

func (res *Result) visit(url *url.URL) error {
	req := newRequest(url)
	var t0, t1, t5 time.Time

	trace := &httptrace.ClientTrace{
		DNSStart: func(_ httptrace.DNSStartInfo) {
			t0 = time.Now()
		},
		DNSDone: func(_ httptrace.DNSDoneInfo) {
			res.dnsLookupTime = float64(time.Since(t0).Milliseconds())
		},
		ConnectStart: func(_, _ string) {
			if t1.IsZero() {
				t1 = time.Now()
			}
		},
		ConnectDone: func(net, addr string, err error) {
			if err != nil {
				log.Fatalf("unable to connect to host %v: %v", addr, err)
			}
			res.tcphandshakeTime = float64(time.Since(t1).Milliseconds())
		},
		TLSHandshakeStart: func() {
			t5 = time.Now()
		},
		TLSHandshakeDone: func(_ tls.ConnectionState, _ error) {
			res.sslhandshakeTime = float64(time.Since(t5).Milliseconds())
		},
	}
	req = req.WithContext(httptrace.WithClientTrace(context.Background(), trace))

	tr := &http.Transport{}

	host, _, err := net.SplitHostPort(req.Host)
	if err != nil {
		host = req.Host
	}

	tr.TLSClientConfig = &tls.Config{
		ServerName:         host,
		InsecureSkipVerify: false,
		MinVersion:         tls.VersionTLS12,
	}

	client := &http.Client{
		Transport: tr,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

func newRequest(url *url.URL) *http.Request {
	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		log.Fatalf("unable to create request: %v", err)
	}
	return req
}

func (s Plugin) FetchMetrics() (map[string]float64, error) {
	res := Result{}
	err := res.visit(s.URL)
	if err != nil {
		return make(map[string]float64), nil
	}
	return map[string]float64{
		"dnsLookupTime":    res.dnsLookupTime,
		"tcphandshakeTime": res.tcphandshakeTime,
		"sslhandshakeTime": res.sslhandshakeTime,
	}, nil
}

func Do() {
	optPrefix := flag.String("prefix", "", "Metric key prefix")
	flag.Parse()
	args := flag.Args()
	url := parseURL(args[0])
	s := Plugin{
		Prefix: *optPrefix,
		URL:    url,
	}
	plugin := mp.NewMackerelPlugin(s)
	plugin.Run()
}
