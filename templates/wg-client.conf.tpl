[Interface]
PrivateKey = {{ .ClientPrivateKey }}
Address = {{ .ClientAddress }}
DNS = {{ .DNS }}

[Peer]
PublicKey = {{ .ServerPublicKey }}
Endpoint = {{ .ServerEndpoint }}
AllowedIPs = {{ .AllowedIPs }}
PersistentKeepalive = 25