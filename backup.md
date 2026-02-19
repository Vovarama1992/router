вот в коде моем

client
dev tun
proto tcp-client
remote basikvpn.ru 443
nobind
resolv-retry infinite
persist-key
persist-tun
remote-cert-tls server
cipher AES-256-GCM
auth SHA256

mssfix 1360
tun-mtu 1400

<tls-crypt>
{{TLS}}
</tls-crypt>

<ca>
{{CA}}
</ca>

<cert>
{{CERT}}
</cert>

<key>
{{KEY}}
</key>

package delivery

import (
"net/http"
"router/internal/domain"
)

type VPNHandler struct {
svc \*domain.Service
}

func NewVPNHandler(svc *domain.Service) *VPNHandler {
return &VPNHandler{svc: svc}
}

// POST /vpn/peer
func (h *VPNHandler) CreatePeer(w http.ResponseWriter, r *http.Request) {
peer, err := h.svc.CreatePeer(r.Context())
if err != nil {
http.Error(w, err.Error(), http.StatusInternalServerError)
return
}

    w.Header().Set("Content-Type", "text/plain")
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(peer.Config))

}

package domain

import (
"context"
"fmt"
"log"
"time"

    "router/internal/infra"
    "router/internal/openvpn"

)

type Peer struct {
Config string
}

type Service struct {
repo \*infra.PeerRepo
}

func NewService(repo *infra.PeerRepo) *Service {
return &Service{repo: repo}
}

func (s *Service) CreatePeer(ctx context.Context) (*Peer, error) {
start := time.Now()

    log.Printf("[domain] CreatePeer start")

    count, err := s.repo.Count(ctx)
    if err != nil {
    	log.Printf("[domain] repo.Count FAILED err=%v", err)
    	return nil, err
    }
    log.Printf("[domain] repo.Count OK count=%d", count)

    name := fmt.Sprintf("peer_%d", count+1)
    log.Printf("[domain] peer name=%s", name)

    client, err := openvpn.CreatePeer(name, "tcp")
    if err != nil {
    	log.Printf("[domain] openvpn.CreatePeer FAILED name=%s err=%v", name, err)
    	return nil, err
    }
    log.Printf("[domain] openvpn.CreatePeer OK name=%s", name)

    if err := s.repo.Save(ctx, count+1, name); err != nil {
    	log.Printf("[domain] repo.Save FAILED id=%d name=%s err=%v", count+1, name, err)
    	return nil, err
    }
    log.Printf("[domain] repo.Save OK id=%d name=%s", count+1, name)

    log.Printf("[domain] CreatePeer done name=%s duration=%s", name, time.Since(start))

    return &Peer{
    	Config: client.Config,
    }, nil

}

package openvpn

import (
"fmt"
"log"
"os"
"os/exec"
"strings"
"time"
)

type Client struct {
Config string
}

const (
pkiDir = "/etc/openvpn/easy-rsa/pki"
caPath = "/etc/openvpn/ca.crt"
tlsPath = "/etc/openvpn/tls-crypt.key"
serverIP = "185.253.8.123"
)

func CreatePeer(name string, transport string) (\*Client, error) {
start := time.Now()

    if transport == "" {
    	transport = "udp"
    }

    keyPath := pkiDir + "/private/" + name + ".key"
    csrPath := pkiDir + "/reqs/" + name + ".csr"
    certPath := pkiDir + "/issued/" + name + ".crt"

    run := func(cmd *exec.Cmd) error {
    	out, err := cmd.CombinedOutput()
    	if err != nil {
    		return fmt.Errorf("cmd error: %v: %s", err, out)
    	}
    	return nil
    }

    if err := run(exec.Command("openssl", "genrsa", "-out", keyPath, "2048")); err != nil {
    	return nil, err
    }

    if err := run(exec.Command(
    	"openssl", "req", "-new",
    	"-key", keyPath,
    	"-out", csrPath,
    	"-subj", "/CN="+name,
    )); err != nil {
    	return nil, err
    }

    if err := run(exec.Command(
    	"openssl", "x509", "-req",
    	"-in", csrPath,
    	"-CA", caPath,
    	"-CAkey", pkiDir+"/private/ca.key",
    	"-CAcreateserial",
    	"-out", certPath,
    	"-days", "825",
    	"-sha256",
    )); err != nil {
    	return nil, err
    }

    read := func(path string) (string, error) {
    	b, err := os.ReadFile(path)
    	if err != nil {
    		return "", err
    	}
    	return string(b), nil
    }

    ca, err := read(caPath)
    if err != nil {
    	return nil, err
    }
    cert, err := read(certPath)
    if err != nil {
    	return nil, err
    }
    key, err := read(keyPath)
    if err != nil {
    	return nil, err
    }

    // FIX tls-crypt
    tlsRaw, err := os.ReadFile(tlsPath)
    if err != nil {
    	return nil, err
    }
    tls := strings.TrimSpace(string(tlsRaw)) + "\n"

    var tplPath string
    if transport == "tcp" {
    	tplPath = "internal/configs/client.conf"
    } else {
    	tplPath = "internal/configs/client_udp.conf"
    }

    tpl, err := read(tplPath)
    if err != nil {
    	return nil, err
    }

    cfg := tpl
    cfg = strings.ReplaceAll(cfg, "{{PEER_NAME}}", name)
    cfg = strings.ReplaceAll(cfg, "{{SERVER_IP}}", serverIP)
    cfg = strings.ReplaceAll(cfg, "{{CA}}", ca)
    cfg = strings.ReplaceAll(cfg, "{{CERT}}", cert)
    cfg = strings.ReplaceAll(cfg, "{{KEY}}", key)
    cfg = strings.ReplaceAll(cfg, "{{TLS}}", tls)

    // FIX PEM newlines
    cfg = strings.ReplaceAll(cfg, "\r\n", "\n")
    cfg = strings.ReplaceAll(cfg, "\n\n\n", "\n\n")

    log.Printf("CreatePeer done %s %s", name, time.Since(start))
    return &Client{Config: cfg}, nil

}

вот конфигии севрера

ubuntu@basov:~$ systemctl status openvpn@server-tcp
systemctl status sslh
systemctl status nginx
● openvpn@server-tcp.service - OpenVPN connection to server-tcp
Loaded: loaded (/lib/systemd/system/openvpn@.service; enabled; vendor pres>
Active: active (running) since Wed 2026-02-18 04:29:56 UTC; 23h ago
Docs: man:openvpn(8)
https://community.openvpn.net/openvpn/wiki/Openvpn24ManPage
https://community.openvpn.net/openvpn/wiki/HOWTO
Main PID: 2026885 (openvpn)
Status: "Initialization Sequence Completed"
Tasks: 1 (limit: 2307)
Memory: 2.7M
CPU: 4.063s
CGroup: /system.slice/system-openvpn.slice/openvpn@server-tcp.service
└─2026885 /usr/sbin/openvpn --daemon ovpn-server-tcp --status /run>

Feb 18 04:29:56 basov systemd[1]: Starting OpenVPN connection to server-tcp...
Feb 18 04:29:56 basov systemd[1]: Started OpenVPN connection to server-tcp.
● sslh.service - SSL/SSH multiplexer
Loaded: loaded (/lib/systemd/system/sslh.service; enabled; vendor preset: >
Active: active (running) since Tue 2026-02-17 06:14:27 UTC; 1 day 22h ago
Docs: man:sslh(8)
Main PID: 1778024 (sslh)
Tasks: 2 (limit: 2307)
Memory: 820.0K
CPU: 20.794s
CGroup: /system.slice/sslh.service
├─1778024 /usr/sbin/sslh --foreground --user sslh --listen 0.0.0.0>
└─1778025 /usr/sbin/sslh --foreground --user sslh --listen 0.0.0.0>

Feb 19 03:26:59 basov sslh[2269591]: tls:connection from 172-236-30-215.ip.lino>
Feb 19 03:27:12 basov sslh[2269636]: tls:connection from 172-236-30-215.ip.lino>
Feb 19 03:43:09 basov sslh[2272266]: openvpn:connection from 176.65.148.19.ptr.>
Feb 19 03:44:37 basov sslh[2272534]: openvpn:connection from 20.169.105.13:5105>
Feb 19 03:44:38 basov sslh[2272536]: openvpn:connection from 20.169.105.13:5106>
Feb 19 03:57:52 basov sslh[2274814]: tls:connection from ns1016175.ip-15-204-18>
Feb 19 03:58:37 basov sslh[2274933]: tls:connection from exit-08.tor.r0cket.net>
Feb 19 03:58:37 basov sslh[2274935]: tls:connection from exit-08.tor.r0cket.net>
Feb 19 04:05:19 basov sslh[2276146]: tls:connection from 43.157.156.190:39162 t>
Feb 19 04:12:50 basov sslh[2277536]: tls:connection from 182.210.203.35.bc.goog>
● nginx.service - A high performance web server and a reverse proxy server
Loaded: loaded (/lib/systemd/system/nginx.service; enabled; vendor preset:>
Active: active (running) since Tue 2026-02-17 06:26:51 UTC; 1 day 21h ago
Docs: man:nginx(8)
Main PID: 1778232 (nginx)
Tasks: 3 (limit: 2307)
Memory: 6.8M
CPU: 7.670s
CGroup: /system.slice/nginx.service
├─1778232 "nginx: master process /usr/sbin/nginx -g daemon on; mas>
├─1778233 "nginx: worker process" "" "" "" "" "" "" "" "" "" "" "">
└─1778234 "nginx: worker process" "" "" "" "" "" "" "" "" "" "" "">

Feb 17 06:26:51 basov systemd[1]: Starting A high performance web server and a >
Feb 17 06:26:51 basov systemd[1]: Started A high performance web server and a r>
ubuntu@basov:~$

ubuntu@basov:~$ sudo ss -lntup
Netid State Recv-Q Send-Q Local Address:Port Peer Address:Port Process  
udp UNCONN 0 0 0.0.0.0:51820 0.0.0.0:_  
udp UNCONN 0 0 127.0.0.53%lo:53 0.0.0.0:_ users:(("systemd-resolve",pid=1733531,fd=13))  
udp UNCONN 0 0 [::]:51820 [::]:_  
tcp LISTEN 0 32 127.0.0.1:1194 0.0.0.0:_ users:(("openvpn",pid=2026885,fd=5))  
tcp LISTEN 0 4096 127.0.0.53%lo:53 0.0.0.0:_ users:(("systemd-resolve",pid=1733531,fd=14))  
tcp LISTEN 0 244 127.0.0.1:5432 0.0.0.0:_ users:(("postgres",pid=1733519,fd=6))  
tcp LISTEN 0 128 0.0.0.0:22 0.0.0.0:_ users:(("sshd",pid=1733490,fd=3))  
tcp LISTEN 0 511 0.0.0.0:80 0.0.0.0:_ users:(("nginx",pid=1778234,fd=7),("nginx",pid=1778233,fd=7),("nginx",pid=1778232,fd=7))
tcp LISTEN 0 511 0.0.0.0:8443 0.0.0.0:_ users:(("nginx",pid=1778234,fd=6),("nginx",pid=1778233,fd=6),("nginx",pid=1778232,fd=6))
tcp LISTEN 0 50 0.0.0.0:443 0.0.0.0:_ users:(("sslh",pid=1778025,fd=3),("sslh",pid=1778024,fd=3))  
tcp LISTEN 0 4096 _:8080 _:_ users:(("router",pid=1734060,fd=8))  
tcp LISTEN 0 128 [::]:22 [::]:_ users:(("sshd",pid=1733490,fd=4))  
tcp LISTEN 0 244 [::1]:5432 [::]:\* users:(("postgres",pid=1733519,fd=5))  
ubuntu@basov:~$

ubuntu@basov:~$ sudo iptables -t nat -L -n -v
Chain PREROUTING (policy ACCEPT 0 packets, 0 bytes)
pkts bytes target prot opt in out source destination  
1295K 66M DOCKER all -- \* \* 0.0.0.0/0 0.0.0.0/0 ADDRTYPE match dst-type LOCAL

Chain INPUT (policy ACCEPT 0 packets, 0 bytes)
pkts bytes target prot opt in out source destination

Chain OUTPUT (policy ACCEPT 0 packets, 0 bytes)
pkts bytes target prot opt in out source destination  
 2 120 DOCKER all -- \* \* 0.0.0.0/0 !127.0.0.0/8 ADDRTYPE match dst-type LOCAL

Chain POSTROUTING (policy ACCEPT 0 packets, 0 bytes)
pkts bytes target prot opt in out source destination  
 0 0 MASQUERADE all -- _ !docker0 172.17.0.0/16 0.0.0.0/0  
 2447 347K MASQUERADE all -- _ eth0 10.8.0.0/24 0.0.0.0/0  
 0 0 MASQUERADE all -- _ eth0 10.8.0.0/24 0.0.0.0/0  
 0 0 MASQUERADE all -- _ eth0 10.8.0.0/24 0.0.0.0/0

Chain DOCKER (2 references)
pkts bytes target prot opt in out source destination  
ubuntu@basov:~$

ubuntu@basov:~$ sudo journalctl -u openvpn@server-tcp -n 100
Feb 06 06:47:12 basov systemd[1]: openvpn@server-tcp.service: Failed with resul>
Feb 06 06:47:17 basov systemd[1]: openvpn@server-tcp.service: Scheduled restart>
Feb 06 06:47:17 basov systemd[1]: Stopped OpenVPN connection to server-tcp.
Feb 06 06:47:17 basov systemd[1]: Starting OpenVPN connection to server-tcp...
Feb 06 06:47:17 basov systemd[1]: Started OpenVPN connection to server-tcp.
Feb 06 06:47:17 basov systemd[1]: openvpn@server-tcp.service: Main process exit>
Feb 06 06:47:17 basov systemd[1]: openvpn@server-tcp.service: Failed with resul>
Feb 06 06:47:23 basov systemd[1]: openvpn@server-tcp.service: Scheduled restart>
Feb 06 06:47:23 basov systemd[1]: Stopped OpenVPN connection to server-tcp.
Feb 06 06:47:23 basov systemd[1]: Starting OpenVPN connection to server-tcp...
Feb 06 06:47:23 basov systemd[1]: Started OpenVPN connection to server-tcp.
Feb 06 06:47:23 basov systemd[1]: openvpn@server-tcp.service: Main process exit>
Feb 06 06:47:23 basov systemd[1]: openvpn@server-tcp.service: Failed with resul>
Feb 06 06:47:28 basov systemd[1]: openvpn@server-tcp.service: Scheduled restart>
Feb 06 06:47:28 basov systemd[1]: Stopped OpenVPN connection to server-tcp.
Feb 06 06:47:28 basov systemd[1]: Starting OpenVPN connection to server-tcp...
Feb 06 06:47:28 basov systemd[1]: Started OpenVPN connection to server-tcp.
Feb 06 06:47:28 basov systemd[1]: openvpn@server-tcp.service: Main process exit>
Feb 06 06:47:28 basov systemd[1]: openvpn@server-tcp.service: Failed with resul>
Feb 06 06:47:33 basov systemd[1]: openvpn@server-tcp.service: Scheduled restart>
Feb 06 06:47:33 basov systemd[1]: Stopped OpenVPN connection to server-tcp.
Feb 06 06:47:33 basov systemd[1]: Starting OpenVPN connection to server-tcp...
Feb 06 06:47:33 basov systemd[1]: Started OpenVPN connection to server-tcp.
ubuntu@basov:~$

ubuntu@basov:~$ sudo journalctl -u sslh -n 100
Feb 18 23:16:16 basov sslh[2226007]: tls:connection from no-reverse-dns-configu>
Feb 18 23:29:46 basov sslh[2228415]: tls:connection from 102.22.20.125:62537 to>
Feb 18 23:37:12 basov sslh[2229694]: tls:connection from 104.164.8.29:21236 to >
Feb 18 23:44:24 basov sslh[2230961]: tls:connection from 101.36.106.134:40642 t>
Feb 18 23:44:26 basov sslh[2230962]: tls:connection from 101.36.106.134:41222 t>
Feb 18 23:44:26 basov sslh[2230963]: tls:connection from 101.36.106.134:41220 t>
Feb 18 23:44:50 basov sslh[2231023]: tls:connection from 165.154.206.250:48614 >
Feb 18 23:44:59 basov sslh[2231060]: tls:connection from 165.154.206.250:51392 >
Feb 18 23:45:01 basov sslh[2231077]: tls:connection from 165.154.206.250:52170 >
Feb 18 23:45:01 basov sslh[2231078]: tls:connection from 165.154.206.250:52172 >
Feb 18 23:49:03 basov sslh[2231789]: tls:connection from form.securityresearch.>
Feb 19 00:10:29 basov sslh[2235449]: tls:connection from igutic.earnningipti.co>
Feb 19 00:13:22 basov sslh[2235902]: tls:connection from 8.211.46.74:46408 to o>
Feb 19 00:13:22 basov sslh[2235903]: tls:connection from 8.211.46.74:46418 to o>
Feb 19 00:13:53 basov sslh[2235980]: tls:connection from 47.251.53.147:25180 to>
Feb 19 00:13:56 basov sslh[2235981]: tls:connection from 47.251.53.147:25186 to>
Feb 19 00:13:58 basov sslh[2235998]: tls:connection from 47.251.53.147:25190 to>
Feb 19 00:14:01 basov sslh[2235999]: tls:connection from 47.251.53.147:28744 to>
Feb 19 00:14:03 basov sslh[2236017]: tls:connection from 47.251.53.147:28750 to>
Feb 19 00:14:06 basov sslh[2236018]: tls:connection from 47.251.53.147:28766 to>
Feb 19 00:14:09 basov sslh[2236030]: tls:connection from 47.251.53.147:28780 to>
Feb 19 00:14:11 basov sslh[2236031]: tls:connection from 47.251.53.147:37766 to>
Feb 19 00:14:14 basov sslh[2236041]: tls:connection from 47.251.53.147:37780 to>
ubuntu@basov:~$

ubuntu@basov:~$ curl ifconfig.me
185.253.8.123ubuntu@basov:~$ ip addr show tun0
336396: tun0: <POINTOPOINT,MULTICAST,NOARP,UP,LOWER_UP> mtu 1500 qdisc fq_codel state UNKNOWN group default qlen 500
link/none
inet 10.8.0.1/24 scope global tun0
valid_lft forever preferred_lft forever
inet6 fe80::265a:a45c:eccf:f5b3/64 scope link stable-privacy
valid_lft forever preferred_lft forever
ubuntu@basov:~$

серверный конфиг

ubuntu@basov:/var/www/router$ cat /etc/openvpn/server-tcp.conf
port 1194
local 127.0.0.1
proto tcp
dev tun

persist-key
persist-tun

topology subnet
server 10.8.0.0 255.255.255.0

push "redirect-gateway def1"
push "dhcp-option DNS 1.1.1.1"
push "dhcp-option DNS 8.8.8.8"

keepalive 10 120

cipher AES-256-GCM
data-ciphers AES-256-GCM
auth SHA256

tun-mtu 1400
mssfix 1360
sndbuf 0
rcvbuf 0

ca /etc/openvpn/ca.crt
cert /etc/openvpn/server.crt
key /etc/openvpn/server.key
dh none

tls-crypt /etc/openvpn/tls-crypt.key

verb 6
log /var/log/openvpn-tcp.log
status /var/log/openvpn-tcp-status.log 5
ubuntu@basov:/var/www/router$
