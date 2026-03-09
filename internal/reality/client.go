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
	realityCfg = "config/reality.json"
)

type Client struct {
	UUID string
	Link string
}

type realityConfig struct {
	Server string   `json:"server"`
	Port   int      `json:"port"`
	SID    string   `json:"sid"`
	SNI    []string `json:"sni"`
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

func loadRealityConfig() (*realityConfig, error) {
	data, err := os.ReadFile(realityCfg)
	if err != nil {
		return nil, err
	}

	var cfg realityConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	if len(cfg.SNI) == 0 {
		return nil, fmt.Errorf("sni list empty")
	}

	return &cfg, nil
}

func randomSNI(list []string) string {
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

	inboundsRaw, ok := cfg["inbounds"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid inbounds")
	}

	for _, ib := range inboundsRaw {
		inbound, ok := ib.(map[string]interface{})
		if !ok {
			continue
		}

		protocol, _ := inbound["protocol"].(string)
		if protocol != "vless" {
			continue
		}

		settings, ok := inbound["settings"].(map[string]interface{})
		if !ok {
			continue
		}

		var clients []interface{}

		if raw := settings["clients"]; raw != nil {
			if arr, ok := raw.([]interface{}); ok {
				clients = arr
			}
		}

		clients = append(clients, map[string]interface{}{
			"id": id,
		})

		settings["clients"] = clients
	}

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

	pbk, err := getPBKFromConfig()
	if err != nil {
		return nil, err
	}

	rcfg, err := loadRealityConfig()
	if err != nil {
		return nil, err
	}

	sni := randomSNI(rcfg.SNI)

	link := fmt.Sprintf(
		"vless://%s@%s:%d?encryption=none&security=reality&sni=%s&fp=chrome&pbk=%s&sid=%s&type=tcp#peer",
		id,
		rcfg.Server,
		rcfg.Port,
		sni,
		pbk,
		rcfg.SID,
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

	rcfg, err := loadRealityConfig()
	if err != nil {
		return "", err
	}

	sni := randomSNI(rcfg.SNI)

	link := fmt.Sprintf(
		"vless://%s@%s:%d?encryption=none&security=reality&sni=%s&fp=chrome&pbk=%s&sid=%s&type=tcp#peer",
		uuid,
		rcfg.Server,
		rcfg.Port,
		sni,
		pbk,
		rcfg.SID,
	)

	return link, nil
}
