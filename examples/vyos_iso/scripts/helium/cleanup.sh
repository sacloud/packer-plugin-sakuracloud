#!/bin/vbash
source /opt/vyatta/etc/functions/script-template

# Clean up
sudo unlink /usr/src/linux
sudo unlink /lib/modules/$(uname -r)/build
sudo aptitude -y remove linux-vyatta-kbuild build-essential
sudo apt-get autoclean
sudo apt-get clean

# Delete Debian squeeze package repository and temporary mirror
delete system package repository community
delete system package repository squeeze
delete system package repository squeeze-lts
set system package repository community components 'main'
set system package repository community distribution 'helium'
set system package repository community url 'http://packages.vyos.net/vyos'
commit
save

# Removing leftover leases and persistent rules
sudo rm -f /var/lib/dhcp3/*

# Removing apt caches
sudo rm -rf /var/cache/apt/*

# Removing hw-id
delete interfaces ethernet eth0 hw-id
commit
save

# Adding a 2 sec delay to the interface up, to make the dhclient happy
echo "pre-up sleep 2" | sudo tee -a /etc/network/interfaces
