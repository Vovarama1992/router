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

	inboundsRaw, ok := cfg["inbounds"].([]interface{})
	if !ok {
		return fmt.Errorf("invalid inbounds")
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

		exists := false
		for _, c := range clients {
			m, ok := c.(map[string]interface{})
			if ok && m["id"] == uuid {
				exists = true
				break
			}
		}

		if !exists {
			clients = append(clients, map[string]interface{}{
				"id": uuid,
			})
		}

		settings["clients"] = clients
	}

	out, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(configPath, out, 0644); err != nil {
		return err
	}

	return exec.Command("systemctl", "restart", "xray").Run()
}
