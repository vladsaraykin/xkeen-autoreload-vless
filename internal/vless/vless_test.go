package vless

import (
	"reflect"
	"testing"
)

const swissLine = "vless://11111111-1111-1111-1111-111111111111@217.79.126.159:443?security=reality&encryption=none&fp=firefox&type=tcp&flow=xtls-rprx-vision&sni=sso27-79.synology.com&pbk=PUBLIC&sid=SHORT#%F0%9F%87%A8%F0%9F%87%AD%20%D0%A6%D1%8E%D1%80%D0%B8%D1%85%2C%20%D0%A8%D0%B2%D0%B5%D0%B9%D1%86%D0%B0%D1%80%D0%B8%D1%8F%2C%20Extra"

func TestParseAndSelect(t *testing.T) {
	n, err := Select([]string{swissLine}, "Швейцария", "")
	if err != nil {
		t.Fatal(err)
	}
	if n.Address != "217.79.126.159" || n.SNI != "sso27-79.synology.com" {
		t.Fatalf("unexpected node: %+v", n)
	}
}

func TestCountries(t *testing.T) {
	netherlands := "vless://11111111-1111-1111-1111-111111111111@1.2.3.4:443?security=reality&sni=a.example.com&pbk=P&sid=S#%F0%9F%87%B3%F0%9F%87%B1%20%D0%90%D0%BC%D1%81%D1%82%D0%B5%D1%80%D0%B4%D0%B0%D0%BC%2C%20%D0%9D%D0%B8%D0%B4%D0%B5%D1%80%D0%BB%D0%B0%D0%BD%D0%B4%D1%8B%2C%20Extra"
	got := Countries([]string{swissLine, netherlands, swissLine})
	want := []string{"Нидерланды", "Швейцария"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Countries() = %#v, want %#v", got, want)
	}
}

func TestLocation(t *testing.T) {
	n, err := Parse(swissLine)
	if err != nil {
		t.Fatal(err)
	}
	city, country, ok := Location(n.Name)
	if !ok || city != "🇨🇭 Цюрих" || country != "Швейцария" {
		t.Fatalf("Location() = %q, %q, %v", city, country, ok)
	}
}
