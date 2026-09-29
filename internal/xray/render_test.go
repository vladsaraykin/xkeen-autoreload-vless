package xray

import (
	"encoding/json"
	"github.com/vladsaraykin/xkeen-autoreload-vless/internal/vless"
	"testing"
)

func TestRenderValidJSON(t *testing.T) {
	b, err := Render(vless.Node{UUID: "u", Address: "1.2.3.4", Port: 443, Security: "reality", Encryption: "none", Fingerprint: "firefox", Network: "tcp", Flow: "xtls-rprx-vision", SNI: "example.com", PublicKey: "p", ShortID: "s"})
	if err != nil {
		t.Fatal(err)
	}
	var v any
	if err := json.Unmarshal(b, &v); err != nil {
		t.Fatal(err)
	}
}
