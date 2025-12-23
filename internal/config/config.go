package config

import "os"

type Config struct {
	// WireGuard
	WGInterface string // wg0
	WGEndpoint  string // <server-ip>:51820
	WGPublicKey string // public key сервера
	WGVpnCIDR   string // 10.0.0.0/24

	// Client config defaults
	ClientDNS string // например 1.1.1.1
	ClientIPs string // обычно 0.0.0.0/0

	// Paths
	ClientTplPath string // templates/wg-client.conf.tpl
}

func Load() *Config {
	return &Config{
		WGInterface: os.Getenv("WG_INTERFACE"),
		WGEndpoint:  os.Getenv("WG_SERVER_ENDPOINT"),
		WGPublicKey: os.Getenv("WG_SERVER_PUBLIC_KEY"),
		WGVpnCIDR:   os.Getenv("WG_SERVER_VPN_CIDR"),

		ClientDNS: os.Getenv("WG_CLIENT_DNS"),
		ClientIPs: os.Getenv("WG_CLIENT_ALLOWED_IPS"),

		ClientTplPath: os.Getenv("WG_CLIENT_TEMPLATE_PATH"),
	}
}
