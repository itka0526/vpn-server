# VPN-Server

Short instructions.

# Install Golang

```bash
wget https://go.dev/dl/go1.21.1.linux-amd64.tar.gz
sudo rm -rf /usr/local/go && tar -C /usr/local -xzf go1.21.1.linux-amd64.tar.gz
echo "export PATH=$PATH:/usr/local/go/bin" >> $HOME/.profile
export PATH=$PATH:/usr/local/go/bin
go version
```

# Configure Golang

```bash
git clone https://github.com/itka0526/vpn-server.git
cd /root/vpn-server/ && go build
vim /root/vpn-server/serverConfig.toml
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
wget -O vpn.sh https://get.vpnsetup.net/ovpn
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
