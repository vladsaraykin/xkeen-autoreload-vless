package vless

import (
	"encoding/json"
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

type Node struct {
	UUID, Address                                    string
	Port                                             int
	Security, Encryption, Fingerprint, Network, Flow string
	SNI, PublicKey, ShortID, Name                    string
	RawOutbound                                      json.RawMessage
}

type Choice struct {
	Key   string
	Label string
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

func Select(nodes []Node, selector, city string) (Node, error) {
	for _, n := range nodes {
		if !matchesSelector(n.Name, selector) {
			continue
		}
		if city != "" && !strings.Contains(strings.ToLower(n.Name), strings.ToLower(city)) {
			continue
		}
		return n, nil
	}
	return Node{}, fmt.Errorf("no VLESS node found for selector %q", selector)
}

func Choices(nodes []Node) []Choice {
	type candidate struct {
		label string
		score int
	}
	grouped := make(map[string]candidate)

	for _, n := range nodes {
		key := selectorKey(n.Name)
		if key == "" {
			continue
		}
		label, score := displayLabel(n.Name)
		if current, ok := grouped[key]; !ok || score > current.score {
			grouped[key] = candidate{label: label, score: score}
		}
	}

	out := make([]Choice, 0, len(grouped))
	for key, c := range grouped {
		out = append(out, Choice{Key: key, Label: c.label})
	}
	sort.Slice(out, func(i, j int) bool {
		return strings.ToLower(out[i].Label) < strings.ToLower(out[j].Label)
	})
	return out
}

func selectorKey(name string) string {
	name = strings.TrimSpace(name)
	if flag := leadingFlag(name); flag != "" {
		return flag
	}
	if _, country, ok := Location(name); ok {
		return strings.ToLower(country)
	}
	part := firstSegment(name)
	if part == "" {
		return ""
	}
	return strings.ToLower(part)
}

func matchesSelector(name, selector string) bool {
	selector = strings.TrimSpace(selector)
	if selector == "" {
		return false
	}
	if leadingFlag(selector) != "" {
		return leadingFlag(name) == leadingFlag(selector)
	}
	return strings.Contains(strings.ToLower(name), strings.ToLower(selector))
}

func displayLabel(name string) (string, int) {
	name = strings.TrimSpace(name)
	flag := leadingFlag(name)
	rest := strings.TrimSpace(strings.TrimPrefix(name, flag))

	if _, country, ok := Location(name); ok && country != "" {
		return strings.TrimSpace(flag + " " + country), 100
	}

	segment := firstSegment(rest)
	lower := strings.ToLower(segment)
	generic := segment == "" ||
		strings.HasPrefix(lower, "lb ") ||
		strings.HasPrefix(lower, "премиум") ||
		strings.HasPrefix(lower, "s6") ||
		strings.HasPrefix(lower, "s8") ||
		strings.HasPrefix(lower, "s12") ||
		strings.HasPrefix(lower, "s13") ||
		strings.HasPrefix(lower, "s14") ||
		strings.HasPrefix(lower, "s15") ||
		strings.HasPrefix(lower, "s16") ||
		strings.HasPrefix(lower, "s17") ||
		strings.HasPrefix(lower, "обход")

	score := 20
	if !generic {
		score = 50
		if hasCyrillic(segment) {
			score = 80
		}
	}
	if segment == "" {
		segment = rest
	}
	if segment == "" {
		segment = name
	}
	return strings.TrimSpace(flag + " " + segment), score
}

func firstSegment(s string) string {
	for _, sep := range []string{" · ", " • ", " [", " ⚡", " 🚫", " ⬇"} {
		if i := strings.Index(s, sep); i >= 0 {
			s = s[:i]
		}
	}
	return strings.TrimSpace(s)
}

func leadingFlag(s string) string {
	r := []rune(strings.TrimSpace(s))
	if len(r) < 2 {
		return ""
	}
	if isRegionalIndicator(r[0]) && isRegionalIndicator(r[1]) {
		return string(r[:2])
	}
	return ""
}

func isRegionalIndicator(r rune) bool {
	return r >= 0x1F1E6 && r <= 0x1F1FF
}

func hasCyrillic(s string) bool {
	for _, r := range s {
		if unicode.In(r, unicode.Cyrillic) {
			return true
		}
	}
	return false
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
