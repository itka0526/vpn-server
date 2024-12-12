sudo apt update
sudo apt install apache2 python3
sudo a2enmod cgi

target_path="/usr/lib/cgi-bin/"
sudo cp -r ./src/cgi/. $target_path

sudo systemctl restart apache2