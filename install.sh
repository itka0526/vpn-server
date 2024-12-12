#!/bin/sh

wget https://go.dev/dl/go1.21.1.linux-amd64.tar.gz
sudo rm -rf /usr/local/go && tar -C /usr/local -xzf go1.21.1.linux-amd64.tar.gz
echo "export PATH=$PATH:/usr/local/go/bin" >> /root/.bashrc

sudo /usr/local/go/bin/go version
sudo /usr/local/go/bin/go build

option=$1
url=""
echo "1. OpenVPN\n2. Wireguard"

case "$option" in
    1)
        echo "You have chosen to install OpenVPN."
        url="https://get.vpnsetup.net/ovpn"
        ;;
    2)
        echo "You have chosen to install WireGuard VPN."
        url="https://get.vpnsetup.net/wg"
        ;;
    *)
        echo "Invalid option. Please choose 1, 2"
        return 1
        ;;
esac

wget -O /root/vpn.sh $url

sudo /root/vpn.sh

sudo crontab -l > /root/cron
echo "@reboot /root/vpn-server/vpn-server >> ~/log.txt 2>&1" >> /root/cron
cat /root/cron
sudo crontab /root/cron

rm /root/cron

/root/vpn-server/vpn-server&
