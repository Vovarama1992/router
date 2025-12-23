[Interface]
PrivateKey = {{ .ServerPrivateKey }}
Address = {{ .ServerAddress }}
ListenPort = {{ .ListenPort }}

PostUp = iptables -A FORWARD -i %i -j ACCEPT; iptables -A FORWARD -o %i -j ACCEPT; iptables -t nat -A POSTROUTING -o {{ .ExternalInterface }} -j MASQUERADE
PostDown = iptables -D FORWARD -i %i -j ACCEPT; iptables -D FORWARD -o %i -j ACCEPT; iptables -t nat -D POSTROUTING -o {{ .ExternalInterface }} -j MASQUERADE