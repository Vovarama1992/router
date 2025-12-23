package infra

import (
	"context"
	"fmt"
	"os/exec"
)

type WGApplier struct {
	Interface string
}

func NewWGApplier(iface string) *WGApplier {
	return &WGApplier{Interface: iface}
}

func (w *WGApplier) ApplyPeer(ctx context.Context, publicKey, address string) error {
	cmd := exec.CommandContext(
		ctx,
		"wg", "set", w.Interface,
		"peer", publicKey,
		"allowed-ips", address,
	)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("wg set failed: %v (%s)", err, string(out))
	}
	return nil
}
