sudo aptitude -y purge open-vm-tools open-vm-modules
cd
sudo mount -o loop linux.iso /mnt
tar zxf /mnt/VMwareTools-*.tar.gz -C /tmp/
sudo /tmp/vmware-tools-distrib/vmware-install.pl -d
rm linux.iso
sudo umount /mnt
