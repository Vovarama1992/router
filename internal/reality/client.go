package reality

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"github.com/google/uuid"
)

const (
	configPath = "/usr/local/etc/xray/config.json"
	server     = "185.253.8.123"
	sni        = "www.cloudflare.com"
	pbk        = "HRFMnOZqUdMPoOeCqlC9uFdPbLTUeGp66AcM-LG0Bd0"
	sid        = "eee842cbf9f8e299"
)

type Client struct {
	UUID string
	Link string
}

type xrayConfig struct {
	Inbounds []struct {
		Settings struct {
			Clients []struct {
				ID string `json:"id"`
			} `json:"clients"`
		} `json:"settings"`
	} `json:"inbounds"`
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

	cmd := exec.Command("systemctl", "restart", "xray")
	if err := cmd.Run(); err != nil {
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
