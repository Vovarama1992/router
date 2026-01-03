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

func CreatePeer(name string) (*Client, error) {
	log.Printf("[vpn] CreatePeer start: name=%s", name)

	cmd := exec.Command(
		"/etc/openvpn/easy-rsa/easyrsa",
		"build-client-full",
		name,
		"nopass",
	)
	cmd.Dir = "/etc/openvpn/easy-rsa"
	cmd.Env = append(os.Environ(), "EASYRSA_BATCH=1")

	if out, err := cmd.CombinedOutput(); err != nil {
		log.Printf("[vpn] easyrsa error (%s): %s", name, out)
		return nil, fmt.Errorf("easy-rsa error: %w", err)
	}
	log.Printf("[vpn] easyrsa OK for %s", name)

	read := func(path string) ([]byte, error) {
		log.Printf("[vpn] read file: %s", path)
		b, err := os.ReadFile(path)
		if err != nil {
			log.Printf("[vpn] read FAILED: %s: %v", path, err)
			return nil, err
		}
		log.Printf("[vpn] read OK: %s (%d bytes)", path, len(b))
		return b, nil
	}

	ca, err := read("/etc/openvpn/ca.crt")
	if err != nil {
		return nil, err
	}

	cert, err := read("/etc/openvpn/easy-rsa/pki/issued/" + name + ".crt")
	if err != nil {
		return nil, err
	}

	key, err := read("/etc/openvpn/easy-rsa/pki/private/" + name + ".key")
	if err != nil {
		return nil, err
	}

	tls, err := read("/etc/openvpn/ta.key")
	if err != nil {
		return nil, err
	}

	tplPath := "internal/configs/client.conf"
	tpl, err := read(tplPath)
	if err != nil {
		return nil, err
	}

	cfg := string(tpl)
	cfg = strings.ReplaceAll(cfg, "{{PEER_NAME}}", name)
	cfg = strings.ReplaceAll(cfg, "{{SERVER_IP}}", "185.253.8.123")
	cfg = strings.ReplaceAll(cfg, "{{CA}}", string(ca))
	cfg = strings.ReplaceAll(cfg, "{{CERT}}", string(cert))
	cfg = strings.ReplaceAll(cfg, "{{KEY}}", string(key))
	cfg = strings.ReplaceAll(cfg, "{{TLS}}", string(tls))

	log.Printf("[vpn] client config built: %d bytes", len(cfg))

	return &Client{Config: cfg}, nil
}
