package subscription

import (
	"context"
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestFetchBase64VLESS(t *testing.T) {
	line := "vless://11111111-1111-1111-1111-111111111111@1.2.3.4:443?security=reality&sni=a.example.com&pbk=P&sid=S#%F0%9F%87%A9%F0%9F%87%AA%20%D0%91%D0%B5%D1%80%D0%BB%D0%B8%D0%BD%2C%20%D0%93%D0%B5%D1%80%D0%BC%D0%B0%D0%BD%D0%B8%D1%8F"
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("User-Agent"); got != DefaultUserAgent {
			t.Errorf("User-Agent = %q, want %q", got, DefaultUserAgent)
		}
		_, _ = w.Write([]byte(base64.StdEncoding.EncodeToString([]byte(line))))
	}))
	defer srv.Close()

	nodes, err := New(5 * time.Second).Fetch(context.Background(), srv.URL)
	if err != nil {
		t.Fatal(err)
	}
	if len(nodes) != 1 || nodes[0].Address != "1.2.3.4" {
		t.Fatalf("unexpected nodes: %+v", nodes)
	}
}

func TestFetchJSONConfigArray(t *testing.T) {
	payload := `[{
	  "outbounds":[{
	    "tag":"proxy",
	    "protocol":"vless",
	    "settings":{"vnext":[{"address":"novosrv.raw-srv.org","port":8444,"users":[{"id":"u","encryption":"none","flow":""}]}]},
	    "streamSettings":{
	      "network":"grpc",
	      "grpcSettings":{"serviceName":"grpc","authority":"","mode":false},
	      "security":"reality",
	      "realitySettings":{"serverName":"novosrv.raw-srv.org","publicKey":"pk","shortId":"sid","spiderX":"/","fingerprint":"firefox"}
	    }
	  }],
	  "remarks":"🇳🇱 Нидерланды · ультра"
	}]`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("User-Agent"); got != "v2rayNG/1.8.19" {
			t.Errorf("User-Agent = %q", got)
		}
		_, _ = w.Write([]byte(payload))
	}))
	defer srv.Close()

	nodes, err := New(5 * time.Second).Fetch(context.Background(), srv.URL)
	if err != nil {
		t.Fatal(err)
	}
	if len(nodes) != 1 {
		t.Fatalf("len(nodes) = %d", len(nodes))
	}
	n := nodes[0]
	if n.Network != "grpc" || n.Address != "novosrv.raw-srv.org" || n.Port != 8444 || n.Name != "🇳🇱 Нидерланды · ультра" {
		t.Fatalf("unexpected node: %+v", n)
	}
	if len(n.RawOutbound) == 0 {
		t.Fatal("RawOutbound is empty")
	}
}
