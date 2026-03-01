package reality

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
)

func AddClient(uuid string) error {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	var cfg map[string]interface{}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return err
	}

	inbounds, ok := cfg["inbounds"].([]interface{})
	if !ok {
		return fmt.Errorf("invalid inbounds")
	}

	var vpnInbound map[string]interface{}

	for _, ib := range inbounds {
		inb := ib.(map[string]interface{})
		if inb["tag"] == "vpn" {
			vpnInbound = inb
			break
		}
	}

	if vpnInbound == nil {
		return fmt.Errorf("vpn inbound not found")
	}

	settings := vpnInbound["settings"].(map[string]interface{})
	clientsRaw, ok := settings["clients"].([]interface{})
	if !ok {
		return fmt.Errorf("clients not found")
	}

	// Проверим, нет ли уже такого клиента
	for _, c := range clientsRaw {
		client := c.(map[string]interface{})
		if client["id"] == uuid {
			return nil // уже есть
		}
	}

	newClient := map[string]interface{}{
		"id":    uuid,
		"level": 0,
		"email": uuid,
	}

	settings["clients"] = append(clientsRaw, newClient)

	out, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(configPath, out, 0644); err != nil {
		return err
	}

	go func() {
		_ = exec.Command("systemctl", "restart", "xray").Run()
	}()

	return nil
}
