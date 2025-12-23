package wg

import (
	"bytes"
	"os/exec"
	"strings"
)

type KeyPair struct {
	Private string
	Public  string
}

func GenerateKeyPair() (*KeyPair, error) {
	// wg genkey
	cmd := exec.Command("wg", "genkey")
	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return nil, err
	}

	privateKey := strings.TrimSpace(out.String())

	// echo <private> | wg pubkey
	cmd = exec.Command("wg", "pubkey")
	cmd.Stdin = strings.NewReader(privateKey)
	out.Reset()
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return nil, err
	}

	publicKey := strings.TrimSpace(out.String())

	return &KeyPair{
		Private: privateKey,
		Public:  publicKey,
	}, nil
}
