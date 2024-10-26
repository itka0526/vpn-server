# VPN-Server

Short instructions.

# Install Golang

```bash
wget "https://go.dev/dl/go1.23.2.linux-amd64.tar.gz"
rm -rf /usr/local/go && tar -C /usr/local -xzf go1.23.2.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
```

# Configure Golang

```bash
git clone https://github.com/itka0526/vpn-server.git
cd /root/vpn-server/ && go build
```

# Install Pihole

```bash
curl -sSL https://install.pi-hole.net | bash
```

# Configure Pihole

```bash
pihole -a -p
vim /etc/lighttpd/lighttpd.conf # port 80 -> 8020
systemctl restart lighttpd.service
```

# Install OVPN (1)

```bash
curl -O https://raw.githubusercontent.com/angristan/openvpn-install/master/openvpn-install.sh
chmod +x vpn.sh
```

# Install WG (2)

```bash
wget -O vpn.sh https://get.vpnsetup.net/wg
chmod +x vpn.sh
```

# Run

```bash
crontab -e
@reboot /root/vpn-server/vpn-server >> ~/log.txt 2>&1
```
