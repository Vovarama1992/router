package reality

import (
	"encoding/json"
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

	inbounds := cfg["inbounds"].([]interface{})
	inb := inbounds[0].(map[string]interface{})

	settings := inb["settings"].(map[string]interface{})
	clients := settings["clients"].([]interface{})

	var newClients []interface{}

	for _, c := range clients {
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

	exec.Command("systemctl", "restart", "xray").Run()

	return nil
}
