package xray

import (
	"encoding/json"

	"github.com/vladsaraykin/xkeen-autoreload-vless/internal/vless"
)

type config struct {
	Outbounds []outbound `json:"outbounds"`
}
type outbound struct {
	Tag            string          `json:"tag"`
	Protocol       string          `json:"protocol"`
	Settings       any             `json:"settings,omitempty"`
	StreamSettings *streamSettings `json:"streamSettings,omitempty"`
}
type streamSettings struct {
	Network         string          `json:"network"`
	Security        string          `json:"security"`
	RealitySettings realitySettings `json:"realitySettings"`
}
type realitySettings struct {
	Fingerprint string `json:"fingerprint"`
	ServerName  string `json:"serverName"`
	PublicKey   string `json:"publicKey"`
	ShortID     string `json:"shortId"`
	SpiderX     string `json:"spiderX"`
}

func Render(n vless.Node) ([]byte, error) {
	c := config{Outbounds: []outbound{
		{
			Tag: "vless-reality", Protocol: "vless",
			Settings: map[string]any{"vnext": []any{map[string]any{
				"address": n.Address, "port": n.Port,
				"users": []any{map[string]any{"id": n.UUID, "flow": n.Flow, "encryption": n.Encryption, "level": 0}},
			}}},
			StreamSettings: &streamSettings{Network: n.Network, Security: n.Security, RealitySettings: realitySettings{
				Fingerprint: n.Fingerprint, ServerName: n.SNI, PublicKey: n.PublicKey, ShortID: n.ShortID, SpiderX: "/",
			}},
		},
		{Tag: "direct", Protocol: "freedom"},
		{Tag: "block", Protocol: "blackhole", Settings: map[string]any{"response": map[string]any{"type": "http"}}},
	}}
	return json.MarshalIndent(c, "", "  ")
}
