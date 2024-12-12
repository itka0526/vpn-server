#!/bin/sh

wget https://go.dev/dl/go1.21.1.linux-amd64.tar.gz
sudo rm -rf /usr/local/go && tar -C /usr/local -xzf go1.21.1.linux-amd64.tar.gz
echo "export PATH=$PATH:/usr/local/go/bin" >> $HOME/.profile
go version

go build

wget -O /root/vpn.sh https://get.vpnsetup.net/wg

sudo crontab -l > /root/cron
echo "@reboot /root/vpn-server/vpn-server >> ~/log.txt 2>&1" >> /root/cron
cat /root/cron
sudo crontab /root/cron
rm /root/cron