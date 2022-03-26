package ssl

import (
	"net/http"
	"net/url"
	"reflect"
	"testing"

	mp "github.com/mackerelio/go-mackerel-plugin"
)

func TestPlugin_GraphDefinition(t *testing.T) {
	u, _ := url.Parse("https://exmaple.com")
	type fields struct {
		Prefix string
		URL    *url.URL
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string]mp.Graphs
	}{
		{
			name: "test",
			fields: fields{
				Prefix: "ssl",
				URL:    u,
			},
			want: map[string]mp.Graphs{
				"": {
					Label: "Ssl_connection_time",
					Unit:  mp.UnitFloat,
					Metrics: []mp.Metrics{
						{Name: "dnsLookupTime", Label: "dnsLookupTime"},
						{Name: "tcphandshakeTime", Label: "tcphandshakeTime"},
						{Name: "sslhandshakeTime", Label: "sslhandshakeTime"},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Plugin{
				Prefix: tt.fields.Prefix,
				URL:    tt.fields.URL,
			}
			if got := s.GraphDefinition(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Plugin.GraphDefinition() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPlugin_MetricKeyPrefix(t *testing.T) {
	u, _ := url.Parse("https://exmaple.com")
	type fields struct {
		Prefix string
		URL    *url.URL
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "ssl_test1",
			fields: fields{
				Prefix: "ssl",
				URL:    u,
			},
			want: "ssl",
		},
		{
			name: "no_prefix",
			fields: fields{
				Prefix: "",
				URL:    u,
			},
			want: "ssl",
		},
		{
			name: "prefix",
			fields: fields{
				Prefix: "test",
				URL:    u,
			},
			want: "test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Plugin{
				Prefix: tt.fields.Prefix,
				URL:    tt.fields.URL,
			}
			if got := s.MetricKeyPrefix(); got != tt.want {
				t.Errorf("Plugin.MetricKeyPrefix() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseURL(t *testing.T) {
	u, _ := url.Parse("https://exmaple.com")

	type args struct {
		uri string
	}
	tests := []struct {
		name string
		args args
		want *url.URL
	}{
		{
			name: "test01",
			args: args{
				uri: "https://exmaple.com",
			},
			want: u,
		},
		{
			name: "test02",
			args: args{
				uri: "exmaple.com",
			},
			want: u,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseURL(tt.args.uri); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResult_visit(t *testing.T) {
	u, _ := url.Parse("https://exmaple.com")
	type fields struct {
		dnsLookupTime    float64
		tcphandshakeTime float64
		sslhandshakeTime float64
	}
	type args struct {
		url *url.URL
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "Test",
			fields: fields{},
			args: args{
				url: u,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := &Result{
				dnsLookupTime:    tt.fields.dnsLookupTime,
				tcphandshakeTime: tt.fields.tcphandshakeTime,
				sslhandshakeTime: tt.fields.sslhandshakeTime,
			}
			if err := res.visit(tt.args.url); (err != nil) != tt.wantErr {
				t.Errorf("Result.visit() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_newRequest(t *testing.T) {
	u, _ := url.Parse("https://exmaple.com")
	req, _ := http.NewRequest("GET", u.String(), nil)
	type args struct {
		url *url.URL
	}
	tests := []struct {
		name string
		args args
		want *http.Request
	}{
		{
			name: "test1",
			args: args{
				url: u,
			},
			want: req,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newRequest(tt.args.url); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPlugin_FetchMetrics(t *testing.T) {
	type fields struct {
		Prefix string
		URL    *url.URL
	}
	tests := []struct {
		name    string
		fields  fields
		want    map[string]float64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Plugin{
				Prefix: tt.fields.Prefix,
				URL:    tt.fields.URL,
			}
			got, err := s.FetchMetrics()
			if (err != nil) != tt.wantErr {
				t.Errorf("Plugin.FetchMetrics() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Plugin.FetchMetrics() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDo(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Do()
		})
	}
}
