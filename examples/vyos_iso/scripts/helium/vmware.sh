sudo aptitude -y purge open-vm-tools
cd
sudo mkdir /mnt/iso
sudo mount -o loop linux.iso /mnt/iso
tar zxf /mnt/iso/VMwareTools-*.tar.gz -C .
sudo ./vmware-tools-distrib/vmware-install.pl -d
sudo umount /mnt/iso
rm linux.iso
rm -rf ./vmware-tools-distrib
