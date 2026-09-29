package vless

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

type Node struct {
	UUID, Address                                    string
	Port                                             int
	Security, Encryption, Fingerprint, Network, Flow string
	SNI, PublicKey, ShortID, Name                    string
}

func Parse(raw string) (Node, error) {
	u, err := url.Parse(raw)
	if err != nil {
		return Node{}, err
	}
	if u.Scheme != "vless" {
		return Node{}, fmt.Errorf("unsupported scheme %q", u.Scheme)
	}
	port, err := strconv.Atoi(u.Port())
	if err != nil {
		return Node{}, fmt.Errorf("invalid port: %w", err)
	}
	q := u.Query()
	n := Node{
		UUID: u.User.Username(), Address: u.Hostname(), Port: port,
		Security:    defaultString(q.Get("security"), "reality"),
		Encryption:  defaultString(q.Get("encryption"), "none"),
		Fingerprint: defaultString(q.Get("fp"), "firefox"),
		Network:     defaultString(q.Get("type"), "tcp"),
		Flow:        defaultString(q.Get("flow"), "xtls-rprx-vision"),
		SNI:         q.Get("sni"), PublicKey: q.Get("pbk"), ShortID: q.Get("sid"),
		Name: u.Fragment,
	}
	if n.UUID == "" || n.Address == "" || n.Port == 0 || n.SNI == "" || n.PublicKey == "" || n.ShortID == "" {
		return Node{}, fmt.Errorf("incomplete VLESS node")
	}
	return n, nil
}

func Select(lines []string, country, city string) (Node, error) {
	country = strings.ToLower(country)
	city = strings.ToLower(city)
	var countryMatch *Node
	for _, line := range lines {
		if !strings.HasPrefix(line, "vless://") {
			continue
		}
		n, err := Parse(line)
		if err != nil {
			continue
		}
		name := strings.ToLower(n.Name)
		if !strings.Contains(name, country) {
			continue
		}
		if city != "" && strings.Contains(name, city) {
			return n, nil
		}
		if countryMatch == nil {
			copy := n
			countryMatch = &copy
		}
	}
	if countryMatch != nil {
		return *countryMatch, nil
	}
	return Node{}, fmt.Errorf("no VLESS node found for country %q", country)
}

func (n Node) StableKey() string {
	return fmt.Sprintf("%s|%s|%d|%s", n.UUID, n.Address, n.Port, n.PublicKey)
}
func defaultString(v, d string) string {
	if v == "" {
		return d
	}
	return v
}
