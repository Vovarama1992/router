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
		Settings struct {
			Clients []struct {
				ID string `json:"id"`
			} `json:"clients"`
		} `json:"settings"`
	} `json:"inbounds"`
}

func parsePublicKey(s string) string {
	for _, line := range strings.Split(s, "\n") {
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "Password:") {
			return strings.TrimSpace(strings.TrimPrefix(line, "Password:"))
		}

		if strings.HasPrefix(line, "PublicKey:") {
			return strings.TrimSpace(strings.TrimPrefix(line, "PublicKey:"))
		}
	}
	return ""
}

func CreateClient() (*Client, error) {
	id := uuid.New().String()

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var cfg xrayConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	priv := cfg.Inbounds[0].StreamSettings.RealitySettings.PrivateKey

	cmd := exec.Command("xray", "x25519", "-i", priv)
	outCmd, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	pub := parsePublicKey(string(outCmd))
	if pub == "" {
		return nil, fmt.Errorf("failed to parse public key")
	}

	cfg.Inbounds[0].Settings.Clients = append(
		cfg.Inbounds[0].Settings.Clients,
		struct {
			ID string `json:"id"`
		}{ID: id},
	)

	out, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return nil, err
	}

	if err := os.WriteFile(configPath, out, 0644); err != nil {
		return nil, err
	}

	cmdRestart := exec.Command("systemctl", "restart", "xray")
	if err := cmdRestart.Run(); err != nil {
		return nil, err
	}

	link := fmt.Sprintf(
		"vless://%s@%s:443?encryption=none&security=reality&sni=%s&fp=chrome&pbk=%s&sid=%s&type=tcp#peer",
		id,
		server,
		sni,
		pub,
		sid,
	)

	return &Client{
		UUID: id,
		Link: link,
	}, nil
}
