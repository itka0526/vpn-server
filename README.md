# VPN-Server

Short instructions. Don't forget to add "\n" to your request { creds: "secret" + "\n"}.

# Build

cd /root/vpn-server/ && go build

# Run

crontab -e
@reboot /root/vpn-server/vpn-server >> ~/log.txt 2>&1
