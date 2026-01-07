package openvpn

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

type Client struct {
	Config string
}

const (
	pkiDir   = "/etc/openvpn/easy-rsa/pki"
	caPath   = "/etc/openvpn/ca.crt"
	tlsPath  = "/etc/openvpn/ta.key"
	serverIP = "185.253.8.123"
)

// transport:
// ""      -> UDP (по умолчанию)
// "udp"   -> UDP
// "tcp"   -> TCP
func CreatePeer(name string, transport string) (*Client, error) {
	if transport == "" {
		transport = "udp"
	}

	log.Printf("[vpn] CreatePeer start: name=%s transport=%s", name, transport)

	keyPath := pkiDir + "/private/" + name + ".key"
	csrPath := pkiDir + "/reqs/" + name + ".csr"
	certPath := pkiDir + "/issued/" + name + ".crt"

	run := func(cmd *exec.Cmd) error {
		out, err := cmd.CombinedOutput()
		if err != nil {
			log.Printf("[vpn] cmd failed: %s\n%s", strings.Join(cmd.Args, " "), out)
			return fmt.Errorf("cmd error: %w", err)
		}
		return nil
	}

	// 1) private key
	if err := run(exec.Command("openssl", "genrsa", "-out", keyPath, "2048")); err != nil {
		return nil, err
	}

	// 2) CSR
	if err := run(exec.Command(
		"openssl", "req", "-new",
		"-key", keyPath,
		"-out", csrPath,
		"-subj", "/CN="+name,
	)); err != nil {
		return nil, err
	}

	// 3) sign cert
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

	read := func(path string) ([]byte, error) {
		return os.ReadFile(path)
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
	tls, err := read(tlsPath)
	if err != nil {
		return nil, err
	}

	var tplPath string
	switch transport {
	case "udp":
		tplPath = "internal/configs/client_udp.conf"
	case "tcp":
		tplPath = "internal/configs/client.conf"
	default:
		return nil, fmt.Errorf("unknown transport: %s", transport)
	}

	tpl, err := read(tplPath)
	if err != nil {
		return nil, err
	}

	cfg := string(tpl)
	cfg = strings.ReplaceAll(cfg, "{{PEER_NAME}}", name)
	cfg = strings.ReplaceAll(cfg, "{{SERVER_IP}}", serverIP)
	cfg = strings.ReplaceAll(cfg, "{{CA}}", string(ca))
	cfg = strings.ReplaceAll(cfg, "{{CERT}}", string(cert))
	cfg = strings.ReplaceAll(cfg, "{{KEY}}", string(key))
	cfg = strings.ReplaceAll(cfg, "{{TLS}}", string(tls))

	log.Printf("[vpn] client config built: %d bytes (%s)", len(cfg), transport)
	return &Client{Config: cfg}, nil
}
