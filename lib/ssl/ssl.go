package ssl

import (
	"crypto/tls"
	"flag"
	"fmt"
	"strings"
	"time"

	mp "github.com/mackerelio/go-mackerel-plugin"
)

type SSLPlugin struct {
	ip         string
	port       int
	serverName string
	Prefix     string
}

type SSLHandShakeResult struct {
	duration   float64
	statusCode int
}

func (s SSLPlugin) GraphDefinition() map[string]mp.Graphs {
	labelPrefix := strings.Title(s.MetricKeyPrefix())
	return map[string]mp.Graphs{
		"": {
			Label: labelPrefix,
			Unit:  mp.UnitFloat,
			Metrics: []mp.Metrics{
				{Name: "seconds", Label: "Seconds"},
			},
		},
	}
}

func (s SSLPlugin) DoSSLHandshake() (SSLHandShakeResult, error) {
	ret := SSLHandShakeResult{}
	cfg := &tls.Config{
		ServerName: s.serverName,
	}
	now := time.Now()
	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", s.ip, s.port), cfg)
	ret.duration = float64(time.Since(now).Milliseconds())
	if err != nil {
		fmt.Println(err)
		return ret, err
	}
	defer conn.Close()
	return ret, nil
}

// SSLハンドシェイクをしてその時間を計測した値を返す関数
func (s SSLPlugin) FetchMetrics() (map[string]float64, error) {
	v, err := s.DoSSLHandshake()
	if err != nil {
		return nil, err
	}
	return map[string]float64{"seconds": v.duration}, nil
}

func (s SSLPlugin) MetricKeyPrefix() string {
	return "SSL"
}

func Do() {
	ip := flag.String("ip", "", "description")
	serverName := flag.String("servername", "", "description")
	port := flag.Int("port", 443, "description")
	flag.Parse()
	s := SSLPlugin{
		Prefix:     "SSL",
		ip:         *ip,
		serverName: *serverName,
		port:       *port,
	}
	plugin := mp.NewMackerelPlugin(s)
	plugin.Run()
}
