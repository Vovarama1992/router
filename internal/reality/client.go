package reality

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/google/uuid"
)

const (
	configPath = "/usr/local/etc/xray/config.json"
	server     = "185.253.8.123"
	sni        = "www.cloudflare.com"
	sid        = "eee842cbf9f8e299"
)

type Client struct {
	UUID string
	Link string
}

type xrayConfig struct {
	Inbounds []struct {
		StreamSettings struct {
			RealitySettings struct {
				PrivateKey string `json:"privateKey"`
			} `json:"realitySettings"`
		} `json:"streamSettings"`
	} `json:"inbounds"`
}

func parsePublicKey(s string) string {
	for _, line := range strings.Split(s, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "PublicKey:") {
			return strings.TrimSpace(strings.TrimPrefix(line, "PublicKey:"))
		}
	}
	return ""
}

func getPBKFromConfig() (string, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return "", err
	}

	var cfg xrayConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return "", err
	}

	priv := cfg.Inbounds[0].StreamSettings.RealitySettings.PrivateKey

	cmd := exec.Command("xray", "x25519", "-i", priv)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	pbk := parsePublicKey(string(out))
	if pbk == "" {
		return "", fmt.Errorf("failed to get public key")
	}

	return pbk, nil
}

func CreateClient() (*Client, error) {
	id := uuid.New().String()

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var cfg map[string]interface{}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	inbounds := cfg["inbounds"].([]interface{})
	inb := inbounds[0].(map[string]interface{})

	settings := inb["settings"].(map[string]interface{})
	clients := settings["clients"].([]interface{})

	clients = append(clients, map[string]interface{}{
		"id": id,
	})

	settings["clients"] = clients

	out, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return nil, err
	}

	if err := os.WriteFile(configPath, out, 0644); err != nil {
		return nil, err
	}

	exec.Command("systemctl", "restart", "xray").Run()

	pbk, err := getPBKFromConfig()
	if err != nil {
		return nil, err
	}

	link := fmt.Sprintf(
		"vless://%s@%s:443?encryption=none&security=reality&sni=%s&fp=chrome&pbk=%s&sid=%s&type=tcp#peer",
		id,
		server,
		sni,
		pbk,
		sid,
	)

	return &Client{
		UUID: id,
		Link: link,
	}, nil
}

func BuildLink(uuid string) (string, error) {
	pbk, err := getPBKFromConfig()
	if err != nil {
		return "", err
	}

	link := fmt.Sprintf(
		"vless://%s@%s:443?encryption=none&security=reality&sni=%s&fp=chrome&pbk=%s&sid=%s&type=tcp#peer",
		uuid,
		server,
		sni,
		pbk,
		sid,
	)

	return link, nil
}
