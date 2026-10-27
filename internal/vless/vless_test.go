package vless

import (
	"encoding/json"
	"testing"
)

const swissLine = "vless://11111111-1111-1111-1111-111111111111@217.79.126.159:443?security=reality&encryption=none&fp=firefox&type=tcp&flow=xtls-rprx-vision&sni=sso27-79.synology.com&pbk=PUBLIC&sid=SHORT#%F0%9F%87%A8%F0%9F%87%AD%20%D0%A6%D1%8E%D1%80%D0%B8%D1%85%2C%20%D0%A8%D0%B2%D0%B5%D0%B9%D1%86%D0%B0%D1%80%D0%B8%D1%8F%2C%20Extra"

func TestParseAndSelect(t *testing.T) {
	n, err := Parse(swissLine)
	if err != nil {
		t.Fatal(err)
	}
	got, err := Select([]Node{n}, "Швейцария", "")
	if err != nil {
		t.Fatal(err)
	}
	if got.Address != "217.79.126.159" || got.SNI != "sso27-79.synology.com" {
		t.Fatalf("unexpected node: %+v", got)
	}
}

func TestChoicesGroupsByFlag(t *testing.T) {
	nodes := []Node{
		{Name: "🇳🇱 Премиум 1 ⚡️ [LTE]"},
		{Name: "🇳🇱 Нидерланды · ультра"},
		{Name: "🇩🇪 LB Germany"},
		{Name: "🇩🇪 Германия · премиум ⚡️"},
	}
	choices := Choices(nodes)
	if len(choices) != 2 {
		t.Fatalf("len(choices) = %d, want 2: %#v", len(choices), choices)
	}
	var foundNL, foundDE bool
	for _, c := range choices {
		switch c.Key {
		case "🇳🇱":
			foundNL = c.Label == "🇳🇱 Нидерланды"
		case "🇩🇪":
			foundDE = c.Label == "🇩🇪 Германия"
		}
	}
	if !foundNL || !foundDE {
		t.Fatalf("unexpected choices: %#v", choices)
	}
}

func TestSelectByFlag(t *testing.T) {
	nodes := []Node{
		{Name: "🇳🇱 Нидерланды · ультра", Address: "nl"},
		{Name: "🇩🇪 Германия · премиум", Address: "de"},
	}
	n, err := Select(nodes, "🇩🇪", "")
	if err != nil {
		t.Fatal(err)
	}
	if n.Address != "de" {
		t.Fatalf("selected address = %q, want de", n.Address)
	}
}

func TestRawOutboundDoesNotAffectStableKey(t *testing.T) {
	n := Node{UUID: "u", Address: "a", Port: 443, PublicKey: "p", RawOutbound: json.RawMessage(`{"tag":"proxy"}`)}
	if got := n.StableKey(); got != "u|a|443|p" {
		t.Fatalf("StableKey() = %q", got)
	}
}
