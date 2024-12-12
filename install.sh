#!/bin/sh

wget https://go.dev/dl/go1.21.1.linux-amd64.tar.gz
sudo rm -rf /usr/local/go && tar -C /usr/local -xzf go1.21.1.linux-amd64.tar.gz
echo "export PATH=$PATH:/usr/local/go/bin" >> /root/.bashrc

sudo /usr/local/go/bin/go version
sudo /usr/local/go/bin/go build

install_vpn() {
    local option=$1
    local url=""
    local script_path="/root/vpn.sh"

    case "$option" in
        1)
            echo "You have chosen to install WireGuard VPN."
            url="https://get.vpnsetup.net/wg"
            ;;
        2)
            echo "You have chosen to install OpenVPN."
            url="https://get.vpnsetup.net/ovpn"
            ;;
        3)
            echo "Exiting the installation script. Goodbye!"
            exit 0
            ;;
        *)
            echo "Invalid option. Please choose 1, 2, or 3."
            return 1
            ;;
    esac

    echo "Downloading the VPN setup script from $url..."
    wget -O "$script_path" "$url"

    if [[ $? -ne 0 ]]; then
        echo "Failed to download the VPN setup script. Please check your internet connection and try again."
        return 1
    fi

    echo "Setting execute permissions for the script..."
    chmod +x "$script_path"

    echo "Executing the VPN setup script..."
    bash "$script_path"

    if [[ $? -eq 0 ]]; then
        echo "VPN installation completed successfully."
    else
        echo "VPN installation failed. Please check the setup script for more details."
        return 1
    fi
}

while true; do
    read -p "Enter your choice [1, 2, or 3]: " user_choice

    # Validate input and proceed with installation
    if [[ "$user_choice" == "1" || "$user_choice" == "2" ]]; then
        install_vpn "$user_choice"
        # If installation was successful, exit the loop
        if [[ $? -eq 0 ]]; then
            break
        else
            echo "An error occurred during installation. Please try again."
        fi
    elif [[ "$user_choice" == "3" ]]; then
        echo "Exiting the installation script. Goodbye!"
        exit 0
    else
        echo "Invalid input. Please enter 1, 2, or 3."
        echo "You can try again."
    fi

    echo ""
done

sudo crontab -l > /root/cron
echo "@reboot /root/vpn-server/vpn-server >> ~/log.txt 2>&1" >> /root/cron
cat /root/cron
sudo crontab /root/cron

rm /root/cron

/root/vpn-server/vpn-server&
