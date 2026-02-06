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

// transport:
// ""      -> UDP (по умолчанию)
// "udp"   -> UDP
// "tcp"   -> TCP
func CreatePeer(name string, transport string) (*Client, error) {
	start := time.Now()

	if transport == "" {
		transport = "udp"
	}

	log.Printf("[openvpn] CreatePeer start name=%s transport=%s", name, transport)

	keyPath := pkiDir + "/private/" + name + ".key"
	csrPath := pkiDir + "/reqs/" + name + ".csr"
	certPath := pkiDir + "/issued/" + name + ".crt"

	run := func(step string, cmd *exec.Cmd) error {
		log.Printf("[openvpn] step=%s cmd=%s", step, strings.Join(cmd.Args, " "))
		out, err := cmd.CombinedOutput()
		if err != nil {
			log.Printf("[openvpn] step=%s FAILED err=%v output=%s", step, err, out)
			return fmt.Errorf("cmd error: %w", err)
		}
		log.Printf("[openvpn] step=%s OK", step)
		return nil
	}

	// 1) private key
	if err := run("genrsa", exec.Command("openssl", "genrsa", "-out", keyPath, "2048")); err != nil {
		return nil, err
	}

	// 2) CSR
	if err := run(
		"csr",
		exec.Command(
			"openssl", "req", "-new",
			"-key", keyPath,
			"-out", csrPath,
			"-subj", "/CN="+name,
		),
	); err != nil {
		return nil, err
	}

	// 3) sign cert
	if err := run(
		"sign-cert",
		exec.Command(
			"openssl", "x509", "-req",
			"-in", csrPath,
			"-CA", caPath,
			"-CAkey", pkiDir+"/private/ca.key",
			"-CAcreateserial",
			"-out", certPath,
			"-days", "825",
			"-sha256",
		),
	); err != nil {
		return nil, err
	}

	read := func(label, path string) ([]byte, error) {
		b, err := os.ReadFile(path)
		if err != nil {
			log.Printf("[openvpn] read %s FAILED path=%s err=%v", label, path, err)
			return nil, err
		}
		log.Printf("[openvpn] read %s OK bytes=%d", label, len(b))
		return b, nil
	}

	ca, err := read("ca", caPath)
	if err != nil {
		return nil, err
	}
	cert, err := read("cert", certPath)
	if err != nil {
		return nil, err
	}
	key, err := read("key", keyPath)
	if err != nil {
		return nil, err
	}
	tls, err := read("tls", tlsPath)
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
		log.Printf("[openvpn] unknown transport=%s", transport)
		return nil, fmt.Errorf("unknown transport: %s", transport)
	}

	log.Printf("[openvpn] using template path=%s", tplPath)

	tpl, err := read("template", tplPath)
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

	log.Printf(
		"[openvpn] CreatePeer done name=%s transport=%s bytes=%d duration=%s",
		name,
		transport,
		len(cfg),
		time.Since(start),
	)

	return &Client{Config: cfg}, nil
}
