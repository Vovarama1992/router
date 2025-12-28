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
	ca, err := os.ReadFile("/etc/openvpn/easy-rsa/pki/ca.crt")
	if err != nil {
		log.Println("openvpn: read ca.crt:", err)
		return nil, err
	}

	cert, err := os.ReadFile("/etc/openvpn/easy-rsa/pki/issued/client.crt")
	if err != nil {
		log.Println("openvpn: read client.crt:", err)
		return nil, err
	}

	key, err := os.ReadFile("/etc/openvpn/easy-rsa/pki/private/client.key")
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
proto udp
remote 185.253.8.123 1194
resolv-retry infinite
nobind
persist-key
persist-tun

remote-cert-tls server
cipher AES-256-GCM
auth SHA256

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

<tls-crypt>
%s
</tls-crypt>
`,
		ca,
		cert,
		key,
		ta,
	)

	log.Println("openvpn: client config generated")

	return &Client{Config: cfg}, nil
}
