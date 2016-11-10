#/bin/sh

# for linux
mkisofs -R -V config-2 -o answer_file.iso iso_files/

# for macOS
# hdiutil makehybrid -iso -joliet -default-volume-name config-2 -o answer_file.iso iso_files/
