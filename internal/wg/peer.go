package wg

import (
	"fmt"

	"router/internal/config"
)

type Peer struct {
	PublicKey  string
	PrivateKey string
	Address    string
	Config     string
}

type CreatePeerParams struct {
	ClientIndex int // 2 -> X.X.X.2
}

func CreatePeer(cfg *config.Config, clientIndex int) (*Peer, error) {
	keys, err := GenerateKeyPair()
	if err != nil {
		return nil, err
	}

	clientIP := fmt.Sprintf("10.0.0.%d/32", clientIndex)

	conf, err := RenderClientConfig(
		cfg.ClientTplPath,
		ClientConfigParams{
			ClientPrivateKey: keys.Private,
			ClientAddress:    clientIP,
			DNS:              cfg.ClientDNS,
			ServerPublicKey:  cfg.WGPublicKey,
			ServerEndpoint:   cfg.WGEndpoint,
			AllowedIPs:       cfg.ClientIPs,
		},
	)
	if err != nil {
		return nil, err
	}

	return &Peer{
		PublicKey:  keys.Public,
		PrivateKey: keys.Private,
		Address:    clientIP,
		Config:     conf,
	}, nil
}
