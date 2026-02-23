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

func derivePublicKey(privateKey string) (string, error) {
	cmd := exec.Command("xray", "x25519", "-i", privateKey)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	var pub string
	lines := string(out)

	for _, line := range splitLines(lines) {
		if len(line) > 10 && line[:10] == "PublicKey:" {
			pub = line[10:]
			pub = trim(pub)
			break
		}
	}

	if pub == "" {
		return "", fmt.Errorf("public key not found")
	}

	return pub, nil
}

func splitLines(s string) []string {
	var lines []string
	current := ""
	for _, r := range s {
		if r == '\n' {
			lines = append(lines, current)
			current = ""
			continue
		}
		current += string(r)
	}
	if current != "" {
		lines = append(lines, current)
	}
	return lines
}

func trim(s string) string {
	for len(s) > 0 && (s[0] == ' ' || s[0] == '\t') {
		s = s[1:]
	}
	for len(s) > 0 && (s[len(s)-1] == ' ' || s[len(s)-1] == '\t') {
		s = s[:len(s)-1]
	}
	return s
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

	pub, err := derivePublicKey(priv)
	if err != nil {
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

	exec.Command("systemctl", "restart", "xray").Run()

	link := buildLink(id, pub)

	return &Client{
		UUID: id,
		Link: link,
	}, nil
}

func buildLink(uuid string, pbk string) string {
	return fmt.Sprintf(
		"vless://%s@%s:443?encryption=none&security=reality&sni=%s&fp=chrome&pbk=%s&sid=%s&type=tcp#peer",
		uuid,
		server,
		sni,
		pbk,
		sid,
	)
}

func BuildLink(uuid string) string {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return ""
	}

	var cfg xrayConfig
	if json.Unmarshal(data, &cfg) != nil {
		return ""
	}

	priv := cfg.Inbounds[0].StreamSettings.RealitySettings.PrivateKey

	pub, err := derivePublicKey(priv)
	if err != nil {
		return ""
	}

	return buildLink(uuid, pub)
}
