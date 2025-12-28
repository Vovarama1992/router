package openvpn

import (
	"fmt"
	"log"
	"os"
)

type Client struct {
	Config string
}

func CreateClient() (*Client, error) {
	ca, err := os.ReadFile("/etc/openvpn/ca.crt")
	if err != nil {
		log.Println("openvpn: read ca.crt:", err)
		return nil, err
	}

	cert, err := os.ReadFile("/etc/openvpn/client.crt")
	if err != nil {
		log.Println("openvpn: read client.crt:", err)
		return nil, err
	}

	key, err := os.ReadFile("/etc/openvpn/client.key")
	if err != nil {
		log.Println("openvpn: read client.key:", err)
		return nil, err
	}

	ta, err := os.ReadFile("/etc/openvpn/ta.key")
	if err != nil {
		log.Println("openvpn: read ta.key:", err)
		return nil, err
	}

	cfg := fmt.Sprintf(`
client
dev tun
proto tcp
remote 185.253.8.123 443

resolv-retry infinite
nobind
persist-key
persist-tun

remote-cert-tls server
cipher AES-256-GCM
auth SHA256
key-direction 1
verb 3

<ca>
%s
</ca>

<cert>
%s
</cert>

<key>
%s
</key>

<tls-auth>
%s
</tls-auth>
`,
		ca,
		cert,
		key,
		ta,
	)

	log.Println("openvpn: client config generated")

	return &Client{Config: cfg}, nil
}
