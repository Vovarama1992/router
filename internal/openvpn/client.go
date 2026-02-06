package openvpn

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Client struct {
	Config string
}

const (
	pkiDir   = "/etc/openvpn/easy-rsa/pki"
	caPath   = "/etc/openvpn/ca.crt"
	tlsPath  = "/etc/openvpn/tls-crypt.key"
	serverIP = "185.253.8.123"
)

func CreatePeer(name string, transport string) (*Client, error) {
	start := time.Now()

	if transport == "" {
		transport = "udp"
	}

	keyPath := pkiDir + "/private/" + name + ".key"
	csrPath := pkiDir + "/reqs/" + name + ".csr"
	certPath := pkiDir + "/issued/" + name + ".crt"

	run := func(cmd *exec.Cmd) error {
		out, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("cmd error: %v: %s", err, out)
		}
		return nil
	}

	if err := run(exec.Command("openssl", "genrsa", "-out", keyPath, "2048")); err != nil {
		return nil, err
	}

	if err := run(exec.Command(
		"openssl", "req", "-new",
		"-key", keyPath,
		"-out", csrPath,
		"-subj", "/CN="+name,
	)); err != nil {
		return nil, err
	}

	if err := run(exec.Command(
		"openssl", "x509", "-req",
		"-in", csrPath,
		"-CA", caPath,
		"-CAkey", pkiDir+"/private/ca.key",
		"-CAcreateserial",
		"-out", certPath,
		"-days", "825",
		"-sha256",
	)); err != nil {
		return nil, err
	}

	read := func(path string) (string, error) {
		b, err := os.ReadFile(path)
		if err != nil {
			return "", err
		}
		return string(b), nil
	}

	ca, err := read(caPath)
	if err != nil {
		return nil, err
	}
	cert, err := read(certPath)
	if err != nil {
		return nil, err
	}
	key, err := read(keyPath)
	if err != nil {
		return nil, err
	}

	// FIX tls-crypt
	tlsRaw, err := os.ReadFile(tlsPath)
	if err != nil {
		return nil, err
	}
	tls := strings.TrimSpace(string(tlsRaw)) + "\n"

	var tplPath string
	if transport == "tcp" {
		tplPath = "internal/configs/client.conf"
	} else {
		tplPath = "internal/configs/client_udp.conf"
	}

	tpl, err := read(tplPath)
	if err != nil {
		return nil, err
	}

	cfg := tpl
	cfg = strings.ReplaceAll(cfg, "{{PEER_NAME}}", name)
	cfg = strings.ReplaceAll(cfg, "{{SERVER_IP}}", serverIP)
	cfg = strings.ReplaceAll(cfg, "{{CA}}", ca)
	cfg = strings.ReplaceAll(cfg, "{{CERT}}", cert)
	cfg = strings.ReplaceAll(cfg, "{{KEY}}", key)
	cfg = strings.ReplaceAll(cfg, "{{TLS}}", tls)

	// FIX PEM newlines
	cfg = strings.ReplaceAll(cfg, "\r\n", "\n")
	cfg = strings.ReplaceAll(cfg, "\n\n\n", "\n\n")

	log.Printf("CreatePeer done %s %s", name, time.Since(start))
	return &Client{Config: cfg}, nil
}
