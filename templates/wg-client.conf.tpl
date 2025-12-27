[Interface]
PrivateKey = {{.ClientPrivateKey}}
Address = {{.ClientAddress}}
DNS = {{.DNS}}
MTU = {{.MTU}}

[Peer]
PublicKey = {{.ServerPublicKey}}
Endpoint = {{.ServerEndpoint}}
AllowedIPs = {{.AllowedIPs}}
PersistentKeepalive = {{.PersistentKeepalive}}