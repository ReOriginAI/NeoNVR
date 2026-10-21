package main

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestHTTPSCertificatePersistenceTrustAndRenewal(t *testing.T) {
	c := defaultConfig()
	c.HTTPSAddr = "127.0.0.1:5443"
	c.HTTPSHosts = []string{"192.168.1.130", "neonvr.local", "fd00::130"}
	path := filepath.Join(t.TempDir(), "neonvr.json")
	config, publicCA, err := prepareHTTPS(c, path)
	if err != nil {
		t.Fatal(err)
	}
	block, remaining := pem.Decode(publicCA)
	if block == nil || block.Type != "CERTIFICATE" || len(remaining) != 0 {
		t.Fatal("public CA download must contain only a certificate")
	}
	roots := x509.NewCertPool()
	if !roots.AppendCertsFromPEM(publicCA) {
		t.Fatal("cannot trust generated CA")
	}
	leaf := config.Certificates[0].Leaf
	for _, name := range append(c.HTTPSHosts, "localhost", "127.0.0.1", "::1") {
		if _, err := leaf.Verify(x509.VerifyOptions{Roots: roots, DNSName: name}); err != nil {
			t.Fatalf("certificate does not verify for %s: %v", name, err)
		}
	}
	for _, name := range []string{"ca.pem", "server.pem"} {
		info, err := os.Stat(filepath.Join(filepath.Dir(path), "tls", name))
		if err != nil || info.Mode().Perm() != 0600 {
			t.Fatalf("private identity permissions: %v %v", info, err)
		}
	}
	restarted, sameCA, err := prepareHTTPS(c, path)
	if err != nil || !bytes.Equal(sameCA, publicCA) || !bytes.Equal(restarted.Certificates[0].Certificate[0], config.Certificates[0].Certificate[0]) {
		t.Fatal("restart changed certificate identity", err)
	}
	c.HTTPSHosts = append(c.HTTPSHosts, "camera-viewer.local")
	updated, sameCA, err := prepareHTTPS(c, path)
	if err != nil || !bytes.Equal(sameCA, publicCA) {
		t.Fatal("adding a hostname changed trusted CA", err)
	}
	if _, err := updated.Certificates[0].Leaf.Verify(x509.VerifyOptions{Roots: roots, DNSName: "camera-viewer.local"}); err != nil {
		t.Fatal(err)
	}

	// A trusted client completes a real TLS handshake without skipping validation.
	server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.TLS == nil {
			t.Error("TLS state missing")
		}
		w.Write([]byte("secure"))
	}))
	server.TLS = updated
	server.StartTLS()
	defer server.Close()
	client := &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{RootCAs: roots}}}
	defer client.CloseIdleConnections()
	response, err := client.Get(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	data, _ := io.ReadAll(response.Body)
	response.Body.Close()
	if string(data) != "secure" {
		t.Fatal("HTTPS request failed")
	}

	// A nearly expired server identity renews under the original CA.
	ca, err := loadTLSIdentity(filepath.Join(filepath.Dir(path), "tls", "ca.pem"))
	if err != nil {
		t.Fatal(err)
	}
	nearExpiry := *updated.Certificates[0].Leaf
	nearExpiry.NotAfter = time.Now().Add(10 * 24 * time.Hour)
	key := updated.Certificates[0].PrivateKey.(crypto.Signer)
	der, err := x509.CreateCertificate(rand.Reader, &nearExpiry, ca.Leaf, key.Public(), ca.PrivateKey.(crypto.Signer))
	if err != nil {
		t.Fatal(err)
	}
	if err := writeTLSIdentity(filepath.Join(filepath.Dir(path), "tls", "server.pem"), der, key); err != nil {
		t.Fatal(err)
	}
	renewed, sameCA, err := prepareHTTPS(c, path)
	if err != nil || !bytes.Equal(sameCA, publicCA) || time.Until(renewed.Certificates[0].Leaf.NotAfter) < 300*24*time.Hour {
		t.Fatal("certificate renewal failed", err)
	}
}

func TestHTTPSDisabledAndCorruptIdentity(t *testing.T) {
	c := defaultConfig()
	path := filepath.Join(t.TempDir(), "neonvr.json")
	if config, ca, err := prepareHTTPS(c, path); config != nil || ca != nil || err != nil {
		t.Fatal("disabled HTTPS should not create certificates")
	}
	if _, err := os.Stat(filepath.Join(filepath.Dir(path), "tls")); !os.IsNotExist(err) {
		t.Fatal("disabled HTTPS created certificate storage")
	}
	c.HTTPSAddr = "127.0.0.1:5443"
	if _, _, err := prepareHTTPS(c, path); err != nil {
		t.Fatal(err)
	}
	caPath := filepath.Join(filepath.Dir(path), "tls", "ca.pem")
	os.WriteFile(caPath, []byte("corrupted"), 0600)
	if _, _, err := prepareHTTPS(c, path); err == nil {
		t.Fatal("corrupt CA was silently replaced, breaking existing trust")
	}
}

func TestHTTPSConfigAndLinks(t *testing.T) {
	c := defaultConfig()
	for _, address := range []string{":0", ":65536", "missing-port", c.Addr, "0.0.0.0:8555"} {
		c.HTTPSAddr = address
		if c.validateHTTPS() == nil {
			t.Fatalf("accepted invalid HTTPS address %q", address)
		}
	}
	c.HTTPSAddr = "0.0.0.0:5443"
	for _, host := range []string{"https://neonvr.local", "192.168.1.130:5443", "", "../private", "bad host", "-bad.local"} {
		c.HTTPSHosts = []string{host}
		if c.validateHTTPS() == nil {
			t.Fatalf("accepted invalid certificate host %q", host)
		}
	}
	for _, tc := range []struct{ host, expected string }{
		{"192.168.1.130:5500", "https://192.168.1.130:5443/streams"},
		{"neonvr.local", "https://neonvr.local:5443/streams"},
		{"[fd00::130]:5500", "https://[fd00::130]:5443/streams"},
	} {
		r := httptest.NewRequest("GET", "http://"+tc.host+"/streams", nil)
		if got := httpsURL(r, c.HTTPSAddr, "/streams"); got != tc.expected {
			t.Fatalf("HTTPS link %q, want %q", got, tc.expected)
		}
	}
}

func TestHTTPSSettingsRestartURLAndOrigin(t *testing.T) {
	for _, disable := range []bool{false, true} {
		c, path, m := fixture(t)
		m.active.HTTPSAddr = "0.0.0.0:5443"
		c.HTTPSAddr = m.active.HTTPSAddr
		if disable {
			c.HTTPSAddr = ""
		}
		raw, _ := os.ReadFile(path)
		data, _ := json.Marshal(map[string]any{"config": c, "revision": configRevision(raw), "restart": true})
		r := httptest.NewRequest("POST", "https://neonvr.test:5443/api/config", bytes.NewReader(data))
		r.Header.Set("Content-Type", "application/json")
		r.Header.Set("X-NeoNVR-Token", m.token)
		r.Header.Set("Origin", "http://neonvr.test:5443")
		denied := httptest.NewRecorder()
		m.handler(denied, r)
		if denied.Code != http.StatusForbidden {
			t.Fatal("HTTPS settings accepted an HTTP origin")
		}
		r.Body = io.NopCloser(bytes.NewReader(data))
		r.Header.Set("Origin", "https://neonvr.test:5443")
		w := httptest.NewRecorder()
		m.handler(w, r)
		var response struct {
			NextURL string `json:"next_url"`
		}
		json.Unmarshal(w.Body.Bytes(), &response)
		want := "https://neonvr.test:5443/settings"
		if disable {
			want = "http://neonvr.test:5500/settings"
		}
		if w.Code != http.StatusOK || response.NextURL != want {
			t.Fatalf("settings restart response: %d %s", w.Code, w.Body.String())
		}
	}
}
