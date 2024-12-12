#!/bin/sh

sudo -s
wget https://go.dev/dl/go1.21.1.linux-amd64.tar.gz
sudo rm -rf /usr/local/go && tar -C /usr/local -xzf go1.21.1.linux-amd64.tar.gz
echo "export PATH=$PATH:/usr/local/go/bin" >> $HOME/.profile
source $HOME/.profile
go version
go build

wget -O $HOME/vpn.sh https://get.vpnsetup.net/wg

sudo crontab -l > $HOME/cron
echo "@reboot $HOME/vpn-server/vpn-server >> ~/log.txt 2>&1" >> $HOME/cron
cat $HOME/cron
sudo crontab $HOME/cron
rm $HOME/cron