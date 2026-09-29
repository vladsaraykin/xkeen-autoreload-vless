package vless

import "testing"

func TestParseAndSelect(t *testing.T) {
	line := "vless://11111111-1111-1111-1111-111111111111@217.79.126.159:443?security=reality&encryption=none&fp=firefox&type=tcp&flow=xtls-rprx-vision&sni=sso27-79.synology.com&pbk=PUBLIC&sid=SHORT#%F0%9F%87%A8%F0%9F%87%AD%20%D0%A6%D1%8E%D1%80%D0%B8%D1%85%2C%20%D0%A8%D0%B2%D0%B5%D0%B9%D1%86%D0%B0%D1%80%D0%B8%D1%8F%2C%20Extra"
	n, err := Select([]string{line}, "Швейцария", "Цюрих")
	if err != nil {
		t.Fatal(err)
	}
	if n.Address != "217.79.126.159" || n.SNI != "sso27-79.synology.com" {
		t.Fatalf("unexpected node: %+v", n)
	}
}
