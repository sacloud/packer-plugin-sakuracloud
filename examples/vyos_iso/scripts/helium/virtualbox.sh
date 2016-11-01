#!/bin/vbash
source /opt/vyatta/etc/functions/script-template

if test -f VBoxGuestAdditions.iso ; then

  # Install dkms for dynamic compiles
  sudo aptitude -y install dkms

  # If libdbus is not installed, virtualbox will not autostart
  sudo aptitude -y install --without-recommends libdbus-1-3

  # Install the VirtualBox guest additions
  sudo mount -o loop VBoxGuestAdditions.iso /mnt
  yes|sudo /bin/sh /mnt/VBoxLinuxAdditions.run --target /tmp/vbox || :
  sudo umount /mnt

  sudo cp -p /tmp/vbox/vboxadd* /etc/init.d/
  sudo rm -rf /tmp/vbox

  # Start the newly build driver
  sudo /etc/init.d/vboxadd start
  sudo insserv vboxadd
  sudo insserv vboxadd-service

  # Clean Up
  rm VBoxGuestAdditions.iso
fi
