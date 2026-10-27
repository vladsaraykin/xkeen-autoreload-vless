package subscription

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/vladsaraykin/xkeen-autoreload-vless/internal/vless"
)

const DefaultUserAgent = "v2rayNG/1.8.19"

type Client struct {
	http      *http.Client
	userAgent string
}

func New(timeout time.Duration) *Client {
	return NewWithOptions(timeout, DefaultUserAgent, false)
}

func NewWithOptions(timeout time.Duration, userAgent string, insecureTLS bool) *Client {
	if userAgent == "" {
		userAgent = DefaultUserAgent
	}
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: insecureTLS} // #nosec G402 -- opt-in compatibility flag.
	return &Client{
		http:      &http.Client{Timeout: timeout, Transport: transport},
		userAgent: userAgent,
	}
}

func (c *Client) Fetch(ctx context.Context, url string) ([]vless.Node, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", c.userAgent)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("subscription request failed: %s", resp.Status)
	}

	raw, err := io.ReadAll(io.LimitReader(resp.Body, 16<<20))
	if err != nil {
		return nil, err
	}
	payload := strings.TrimSpace(string(raw))
	if payload == "" {
		return nil, fmt.Errorf("subscription response is empty")
	}

	if strings.HasPrefix(payload, "[") || strings.HasPrefix(payload, "{") {
		if nodes, err := parseJSON([]byte(payload)); err == nil && len(nodes) > 0 {
			return nodes, nil
		}
	}

	decoded, err := decodeBase64(payload)
	if err != nil {
		return nil, fmt.Errorf("unsupported subscription format: expected Base64 VLESS list or JSON config array: %w", err)
	}
	return parseVLESSLines(string(decoded))
}

func parseVLESSLines(payload string) ([]vless.Node, error) {
	var nodes []vless.Node
	for _, line := range strings.Split(payload, "\n") {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "vless://") {
			continue
		}
		node, err := vless.Parse(line)
		if err == nil {
			nodes = append(nodes, node)
		}
	}
	if len(nodes) == 0 {
		return nil, fmt.Errorf("no valid VLESS nodes found")
	}
	return nodes, nil
}

type jsonConfig struct {
	Outbounds []json.RawMessage `json:"outbounds"`
	Remarks   string            `json:"remarks"`
}

type jsonOutbound struct {
	Tag      string `json:"tag"`
	Protocol string `json:"protocol"`
	Settings struct {
		VNext []struct {
			Address string `json:"address"`
			Port    int    `json:"port"`
			Users   []struct {
				ID         string `json:"id"`
				Encryption string `json:"encryption"`
				Flow       string `json:"flow"`
			} `json:"users"`
		} `json:"vnext"`
	} `json:"settings"`
	StreamSettings struct {
		Network         string `json:"network"`
		Security        string `json:"security"`
		RealitySettings struct {
			ServerName  string `json:"serverName"`
			PublicKey   string `json:"publicKey"`
			ShortID     string `json:"shortId"`
			Fingerprint string `json:"fingerprint"`
		} `json:"realitySettings"`
	} `json:"streamSettings"`
}

func parseJSON(raw []byte) ([]vless.Node, error) {
	var configs []jsonConfig
	if err := json.Unmarshal(raw, &configs); err != nil {
		var single jsonConfig
		if errSingle := json.Unmarshal(raw, &single); errSingle != nil {
			return nil, err
		}
		configs = []jsonConfig{single}
	}

	var nodes []vless.Node
	for _, cfg := range configs {
		for _, rawOutbound := range cfg.Outbounds {
			var outbound jsonOutbound
			if err := json.Unmarshal(rawOutbound, &outbound); err != nil || outbound.Protocol != "vless" {
				continue
			}
			if len(outbound.Settings.VNext) == 0 || len(outbound.Settings.VNext[0].Users) == 0 {
				continue
			}
			vnext := outbound.Settings.VNext[0]
			user := vnext.Users[0]
			reality := outbound.StreamSettings.RealitySettings
			if user.ID == "" || vnext.Address == "" || vnext.Port == 0 || reality.PublicKey == "" {
				continue
			}
			nodes = append(nodes, vless.Node{
				UUID:        user.ID,
				Address:     vnext.Address,
				Port:        vnext.Port,
				Encryption:  defaultString(user.Encryption, "none"),
				Flow:        user.Flow,
				Network:     defaultString(outbound.StreamSettings.Network, "tcp"),
				Security:    defaultString(outbound.StreamSettings.Security, "reality"),
				SNI:         reality.ServerName,
				PublicKey:   reality.PublicKey,
				ShortID:     reality.ShortID,
				Fingerprint: defaultString(reality.Fingerprint, "firefox"),
				Name:        cfg.Remarks,
				RawOutbound: append(json.RawMessage(nil), rawOutbound...),
			})
		}
	}
	if len(nodes) == 0 {
		return nil, fmt.Errorf("JSON subscription contains no supported VLESS outbounds")
	}
	return nodes, nil
}

func decodeBase64(s string) ([]byte, error) {
	variants := []*base64.Encoding{base64.StdEncoding, base64.RawStdEncoding, base64.URLEncoding, base64.RawURLEncoding}
	for _, enc := range variants {
		if b, err := enc.DecodeString(s); err == nil {
			return b, nil
		}
	}
	return nil, fmt.Errorf("unsupported base64 payload")
}

func defaultString(v, fallback string) string {
	if v == "" {
		return fallback
	}
	return v
}
