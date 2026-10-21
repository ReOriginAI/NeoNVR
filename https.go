package main

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func (c Config) validateHTTPS() error {
	if c.HTTPSAddr != "" {
		host, port, err := net.SplitHostPort(c.HTTPSAddr)
		if err != nil {
			return fmt.Errorf("https_addr: %w", err)
		}
		n, err := strconv.Atoi(port)
		if err != nil || n < 1 || n > 65535 {
			return errors.New("https_addr must use a port between 1 and 65535")
		}
		_, httpPort, _ := net.SplitHostPort(c.Addr)
		httpPortNumber, _ := strconv.Atoi(httpPort)
		if n == c.LiveWebRTCPort || n == httpPortNumber {
			return errors.New("HTTPS, HTTP and WebRTC must use different ports")
		}
		if _, err := net.ResolveTCPAddr("tcp", net.JoinHostPort(host, port)); err != nil {
			return fmt.Errorf("https_addr: %w", err)
		}
	}
	for _, host := range c.HTTPSHosts {
		if net.ParseIP(host) != nil {
			continue
		}
		if len(host) == 0 || len(host) > 253 || strings.TrimSpace(host) != host {
			return errors.New("https_hosts must contain hostnames or IP addresses without ports or URL prefixes")
		}
		for _, label := range strings.Split(host, ".") {
			if len(label) == 0 || len(label) > 63 || label[0] == '-' || label[len(label)-1] == '-' {
				return fmt.Errorf("invalid HTTPS hostname: %s", host)
			}
			for _, char := range label {
				if !(char >= 'a' && char <= 'z' || char >= 'A' && char <= 'Z' || char >= '0' && char <= '9' || char == '-') {
					return fmt.Errorf("invalid HTTPS hostname: %s", host)
				}
			}
		}
	}
	return nil
}

func httpsNames(c Config) []string {
	names := append([]string{"localhost", "127.0.0.1", "::1", "neonvr"}, c.HTTPSHosts...)
	if hostname, err := os.Hostname(); err == nil && hostname != "" {
		names = append(names, hostname)
	}
	for _, address := range []string{c.Addr, c.HTTPSAddr} {
		host, _, _ := net.SplitHostPort(address)
		if host != "" && host != "0.0.0.0" && host != "::" {
			names = append(names, host)
		}
	}
	if addresses, err := net.InterfaceAddrs(); err == nil {
		for _, address := range addresses {
			if ip, _, err := net.ParseCIDR(address.String()); err == nil && !ip.IsUnspecified() {
				names = append(names, ip.String())
			}
		}
	}
	return names
}

// Store each certificate and its private key in one atomic file. The config
// directory is mounted persistently, so trusting the CA survives image rebuilds.
func writeTLSIdentity(path string, cert []byte, key crypto.Signer) error {
	keyDER, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		return err
	}
	data := append(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: cert}), pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: keyDER})...)
	f, err := os.CreateTemp(filepath.Dir(path), ".tls-*")
	if err != nil {
		return err
	}
	defer os.Remove(f.Name())
	defer f.Close()
	if _, err = f.Write(data); err != nil {
		return err
	}
	if err = f.Sync(); err != nil {
		return err
	}
	if err = f.Close(); err != nil {
		return err
	}
	return os.Rename(f.Name(), path)
}

func loadTLSIdentity(path string) (tls.Certificate, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return tls.Certificate{}, err
	}
	cert, err := tls.X509KeyPair(data, data)
	if err != nil {
		return cert, err
	}
	cert.Leaf, err = x509.ParseCertificate(cert.Certificate[0])
	return cert, err
}

func certificateTemplate(name string, lifetime time.Duration) (*x509.Certificate, error) {
	serial, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return nil, err
	}
	return &x509.Certificate{SerialNumber: serial, Subject: pkix.Name{CommonName: name}, NotBefore: time.Now().Add(-time.Hour), NotAfter: time.Now().Add(lifetime), BasicConstraintsValid: true, KeyUsage: x509.KeyUsageDigitalSignature}, nil
}

func prepareHTTPS(c Config, configPath string) (*tls.Config, []byte, error) {
	if c.HTTPSAddr == "" {
		return nil, nil, nil
	}
	dir := filepath.Join(filepath.Dir(configPath), "tls")
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, nil, err
	}
	caPath, serverPath := filepath.Join(dir, "ca.pem"), filepath.Join(dir, "server.pem")
	ca, err := loadTLSIdentity(caPath)
	if os.IsNotExist(err) {
		template, templateErr := certificateTemplate("NeoNVR local CA", 10*365*24*time.Hour)
		if templateErr != nil {
			return nil, nil, templateErr
		}
		template.IsCA = true
		template.KeyUsage |= x509.KeyUsageCertSign
		key, keyErr := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if keyErr != nil {
			return nil, nil, keyErr
		}
		der, createErr := x509.CreateCertificate(rand.Reader, template, template, key.Public(), key)
		if createErr != nil {
			return nil, nil, createErr
		}
		if err = writeTLSIdentity(caPath, der, key); err != nil {
			return nil, nil, err
		}
		ca, err = loadTLSIdentity(caPath)
	}
	if err != nil {
		return nil, nil, fmt.Errorf("load HTTPS CA: %w", err)
	}
	if !ca.Leaf.IsCA || time.Now().Before(ca.Leaf.NotBefore) || time.Now().After(ca.Leaf.NotAfter) {
		return nil, nil, errors.New("HTTPS CA is invalid or expired; replace config/tls and trust the new certificate")
	}
	roots := x509.NewCertPool()
	roots.AddCert(ca.Leaf)
	server, err := loadTLSIdentity(serverPath)
	if err != nil && !os.IsNotExist(err) {
		return nil, nil, fmt.Errorf("load HTTPS server certificate: %w", err)
	}
	names := httpsNames(c)
	usable := err == nil
	if usable {
		_, err = server.Leaf.Verify(x509.VerifyOptions{Roots: roots, CurrentTime: time.Now().Add(30 * 24 * time.Hour)})
		usable = err == nil
		for _, name := range names {
			usable = usable && server.Leaf.VerifyHostname(name) == nil
		}
	}
	if !usable {
		template, err := certificateTemplate("NeoNVR", 365*24*time.Hour)
		if err != nil {
			return nil, nil, err
		}
		template.ExtKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth}
		if template.NotAfter.After(ca.Leaf.NotAfter) {
			template.NotAfter = ca.Leaf.NotAfter
		}
		seen := map[string]bool{}
		for _, name := range names {
			if seen[name] {
				continue
			}
			seen[name] = true
			if ip := net.ParseIP(name); ip != nil {
				template.IPAddresses = append(template.IPAddresses, ip)
			} else {
				template.DNSNames = append(template.DNSNames, name)
			}
		}
		key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			return nil, nil, err
		}
		der, err := x509.CreateCertificate(rand.Reader, template, ca.Leaf, key.Public(), ca.PrivateKey.(crypto.Signer))
		if err != nil {
			return nil, nil, err
		}
		if err := writeTLSIdentity(serverPath, der, key); err != nil {
			return nil, nil, err
		}
		server, err = loadTLSIdentity(serverPath)
		if err != nil {
			return nil, nil, err
		}
	}
	publicCA := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: ca.Certificate[0]})
	return &tls.Config{MinVersion: tls.VersionTLS12, Certificates: []tls.Certificate{server}}, publicCA, nil
}

// Keep HTTP available for recordings and certificate setup. This explicit link
// preserves the hostname the browser can reach, including bracketed IPv6.
func httpsURL(r *http.Request, address, path string) string {
	if address == "" {
		return ""
	}
	_, port, _ := net.SplitHostPort(address)
	u := &url.URL{Scheme: "http", Host: r.Host}
	return (&url.URL{Scheme: "https", Host: net.JoinHostPort(u.Hostname(), port), Path: path}).String()
}
