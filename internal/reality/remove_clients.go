package reality

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
)

func RemoveClients(uuids []string) error {
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

	var newClients []interface{}

	for _, c := range clientsRaw {
		client := c.(map[string]interface{})
		id := client["id"].(string)

		keep := true
		for _, u := range uuids {
			if id == u {
				keep = false
				break
			}
		}

		if keep {
			newClients = append(newClients, client)
		}
	}

	settings["clients"] = newClients

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
