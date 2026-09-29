package vless

import (
	"fmt"
	"net/url"
	"sort"
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
		_, nodeCountry, ok := Location(n.Name)
		if !ok || !strings.EqualFold(nodeCountry, country) {
			continue
		}
		if city != "" {
			nodeCity, _, _ := Location(n.Name)
			if strings.EqualFold(nodeCity, city) {
				return n, nil
			}
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

func Countries(lines []string) []string {
	unique := make(map[string]string)
	for _, line := range lines {
		if !strings.HasPrefix(line, "vless://") {
			continue
		}
		n, err := Parse(line)
		if err != nil {
			continue
		}
		_, country, ok := Location(n.Name)
		if !ok {
			continue
		}
		key := strings.ToLower(country)
		if _, exists := unique[key]; !exists {
			unique[key] = country
		}
	}
	out := make([]string, 0, len(unique))
	for _, country := range unique {
		out = append(out, country)
	}
	sort.Slice(out, func(i, j int) bool {
		return strings.ToLower(out[i]) < strings.ToLower(out[j])
	})
	return out
}

func Location(name string) (city, country string, ok bool) {
	parts := strings.Split(name, ",")
	if len(parts) < 2 {
		return "", "", false
	}
	city = strings.TrimSpace(parts[0])
	country = strings.TrimSpace(parts[1])
	if city == "" || country == "" {
		return "", "", false
	}
	return city, country, true
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
