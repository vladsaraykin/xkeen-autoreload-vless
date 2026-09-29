package setup

import (
	"context"
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/vladsaraykin/xkeen-autoreload-vless/internal/config"
)

func TestRun(t *testing.T) {
	line := "vless://11111111-1111-1111-1111-111111111111@217.79.126.159:443?security=reality&sni=a.example.com&pbk=P&sid=S#%F0%9F%87%A8%F0%9F%87%AD%20%D0%A6%D1%8E%D1%80%D0%B8%D1%85%2C%20%D0%A8%D0%B2%D0%B5%D0%B9%D1%86%D0%B0%D1%80%D0%B8%D1%8F%2C%20Extra"
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(base64.StdEncoding.EncodeToString([]byte(line))))
	}))
	defer srv.Close()

	path := filepath.Join(t.TempDir(), "config.env")
	cfg := config.Config{
		ConfigPath: path,
		OutputPath: "/tmp/out.json",
		StatePath: "/tmp/state.json",
		XKeenCommand: "xkeen",
		Interval: time.Hour,
		HTTPTimeout: 5 * time.Second,
	}
	input := strings.NewReader(srv.URL + "\n1\n")
	var output strings.Builder
	got, err := Run(context.Background(), cfg, input, &output)
	if err != nil {
		t.Fatal(err)
	}
	if got.Country != "Швейцария" || got.SubscriptionURL != srv.URL {
		t.Fatalf("unexpected config: %+v", got)
	}
	if !strings.Contains(output.String(), "Швейцария") {
		t.Fatalf("country list missing from output: %s", output.String())
	}
}
