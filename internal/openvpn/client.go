package openvpn

import (
	"fmt"
	"os"
)

type Client struct {
	Config string
}

func CreateClient() (*Client, error) {
	ca, err := os.ReadFile("/etc/openvpn/easy-rsa/pki/ca.crt")
	if err != nil {
		return nil, err
	}

	cert, err := os.ReadFile("/etc/openvpn/easy-rsa/pki/issued/client.crt")
	if err != nil {
		return nil, err
	}

	key, err := os.ReadFile("/etc/openvpn/easy-rsa/pki/private/client.key")
	if err != nil {
		return nil, err
	}

	ta, err := os.ReadFile("/etc/openvpn/easy-rsa/ta.key")
	if err != nil {
		return nil, err
	}

	cfg := fmt.Sprintf(`
client
dev tun
proto tcp
remote YOUR_SERVER_IP 443
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
		ca, cert, key, ta,
	)

	return &Client{Config: cfg}, nil
}
