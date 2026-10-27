package xray

import (
	"encoding/json"
	"testing"

	"github.com/vladsaraykin/xkeen-autoreload-vless/internal/vless"
)

func TestRenderValidJSON(t *testing.T) {
	b, err := Render(vless.Node{
		UUID: "u", Address: "1.2.3.4", Port: 443,
		Security: "reality", Encryption: "none", Fingerprint: "firefox",
		Network: "tcp", Flow: "xtls-rprx-vision",
		SNI: "example.com", PublicKey: "p", ShortID: "s",
	})
	if err != nil {
		t.Fatal(err)
	}
	var v any
	if err := json.Unmarshal(b, &v); err != nil {
		t.Fatal(err)
	}
}

func TestRenderPreservesRawJSONOutbound(t *testing.T) {
	raw := json.RawMessage(`{
	  "tag":"proxy",
	  "protocol":"vless",
	  "settings":{"vnext":[{"address":"host","port":8444,"users":[{"id":"u","encryption":"none","flow":""}]}]},
	  "streamSettings":{
	    "network":"grpc",
	    "grpcSettings":{"serviceName":"grpc","authority":"","mode":false},
	    "security":"reality",
	    "realitySettings":{"serverName":"sni","publicKey":"pk","shortId":"sid","fingerprint":"firefox"}
	  },
	  "mux":{"enabled":true,"concurrency":8}
	}`)
	b, err := Render(vless.Node{RawOutbound: raw})
	if err != nil {
		t.Fatal(err)
	}
	var cfg struct {
		Outbounds []map[string]any `json:"outbounds"`
	}
	if err := json.Unmarshal(b, &cfg); err != nil {
		t.Fatal(err)
	}
	if len(cfg.Outbounds) != 3 {
		t.Fatalf("outbounds = %d, want 3", len(cfg.Outbounds))
	}
	proxy := cfg.Outbounds[0]
	if proxy["tag"] != "vless-reality" {
		t.Fatalf("tag = %#v", proxy["tag"])
	}
	stream := proxy["streamSettings"].(map[string]any)
	if stream["network"] != "grpc" {
		t.Fatalf("network = %#v", stream["network"])
	}
	if _, ok := stream["grpcSettings"]; !ok {
		t.Fatal("grpcSettings was not preserved")
	}
	if _, ok := proxy["mux"]; !ok {
		t.Fatal("mux was not preserved")
	}
}
