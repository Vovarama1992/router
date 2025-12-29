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
	cmd := exec.Command(
		"/etc/openvpn/easy-rsa/easyrsa",
		"build-client-full",
		name,
		"nopass",
	)
	cmd.Dir = "/etc/openvpn/easy-rsa"
	cmd.Env = append(os.Environ(), "EASYRSA_BATCH=1")

	if out, err := cmd.CombinedOutput(); err != nil {
		log.Printf("[openvpn] easy-rsa failed (%s): %s", name, out)
		return nil, fmt.Errorf("easy-rsa error: %w", err)
	}

	ca, err := os.ReadFile("/etc/openvpn/ca.crt")
	if err != nil {
		log.Printf("[openvpn] read ca.crt failed: %v", err)
		return nil, err
	}

	cert, err := os.ReadFile("/etc/openvpn/easy-rsa/pki/issued/" + name + ".crt")
	if err != nil {
		log.Printf("[openvpn] read cert failed (%s): %v", name, err)
		return nil, err
	}

	key, err := os.ReadFile("/etc/openvpn/easy-rsa/pki/private/" + name + ".key")
	if err != nil {
		log.Printf("[openvpn] read key failed (%s): %v", name, err)
		return nil, err
	}

	tls, err := os.ReadFile("/etc/openvpn/ta.key")
	if err != nil {
		log.Printf("[openvpn] read tls key failed: %v", err)
		return nil, err
	}

	tpl, err := os.ReadFile("internal/configs/client.conf")
	if err != nil {
		log.Printf("[openvpn] read template failed: %v", err)
		return nil, err
	}

	cfg := string(tpl)
	cfg = strings.ReplaceAll(cfg, "{{PEER_NAME}}", name)
	cfg = strings.ReplaceAll(cfg, "{{SERVER_IP}}", "185.253.8.123")
	cfg = strings.ReplaceAll(cfg, "{{CA}}", string(ca))
	cfg = strings.ReplaceAll(cfg, "{{CERT}}", string(cert))
	cfg = strings.ReplaceAll(cfg, "{{KEY}}", string(key))
	cfg = strings.ReplaceAll(cfg, "{{TLS}}", string(tls))

	return &Client{Config: cfg}, nil
}
