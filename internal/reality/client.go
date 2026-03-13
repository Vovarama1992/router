package reality

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	configPath = "/usr/local/etc/xray/config.json"
	serverAddr = "straightfunctor.xyz"
)

type Client struct {
	UUID string
	Link string
}

type xrayConfig struct {
	Inbounds []struct {
		Port     int `json:"port"`
		Settings struct {
			Clients []struct {
				ID string `json:"id"`
			} `json:"clients"`
		} `json:"settings"`
		StreamSettings struct {
			RealitySettings struct {
				PrivateKey  string   `json:"privateKey"`
				ShortIds    []string `json:"shortIds"`
				ServerNames []string `json:"serverNames"`
			} `json:"realitySettings"`
		} `json:"streamSettings"`
	} `json:"inbounds"`
}

func random(list []string) string {
	rand.Seed(time.Now().UnixNano())
	return list[rand.Intn(len(list))]
}

func parsePublicKey(s string) string {
	for _, line := range strings.Split(s, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "Password:") {
			return strings.TrimSpace(strings.TrimPrefix(line, "Password:"))
		}
	}
	return ""
}

func loadConfig() (*xrayConfig, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var cfg xrayConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	if len(cfg.Inbounds) == 0 {
		return nil, fmt.Errorf("no inbounds")
	}

	return &cfg, nil
}

func getPBK(priv string) (string, error) {
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

	cfg, err := loadConfig()
	if err != nil {
		return nil, err
	}

	in := &cfg.Inbounds[0]

	id := uuid.New().String()

	in.Settings.Clients = append(in.Settings.Clients, struct {
		ID string `json:"id"`
	}{
		ID: id,
	})

	out, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return nil, err
	}

	if err := os.WriteFile(configPath, out, 0644); err != nil {
		return nil, err
	}

	if err := exec.Command("systemctl", "restart", "xray").Run(); err != nil {
		return nil, err
	}

	pbk, err := getPBK(in.StreamSettings.RealitySettings.PrivateKey)
	if err != nil {
		return nil, err
	}

	sni := random(in.StreamSettings.RealitySettings.ServerNames)
	sid := random(in.StreamSettings.RealitySettings.ShortIds)

	link := fmt.Sprintf(
		"vless://%s@%s:%d?encryption=none&security=reality&sni=%s&fp=chrome&pbk=%s&sid=%s&type=tcp#peer",
		id,
		serverAddr,
		in.Port,
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

	cfg, err := loadConfig()
	if err != nil {
		return "", err
	}

	in := &cfg.Inbounds[0]

	pbk, err := getPBK(in.StreamSettings.RealitySettings.PrivateKey)
	if err != nil {
		return "", err
	}

	sni := random(in.StreamSettings.RealitySettings.ServerNames)
	sid := random(in.StreamSettings.RealitySettings.ShortIds)

	link := fmt.Sprintf(
		"vless://%s@%s:%d?encryption=none&security=reality&sni=%s&fp=chrome&pbk=%s&sid=%s&type=tcp#peer",
		uuid,
		serverAddr,
		in.Port,
		sni,
		pbk,
		sid,
	)

	return link, nil
}
