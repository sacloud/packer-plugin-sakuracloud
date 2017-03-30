#!/bin/bash

set -e
set -x

VERSION="$1"
SRC_URL="https://github.com/sacloud/packer-builder-sakuracloud/releases/download/v${VERSION}/packer-builder-sakuracloud_darwin-amd64.zip"
SHA256_SRC=`curl -sL -o - "$SRC_URL" | openssl dgst -sha256`

# clone
git clone --depth=50 --branch=master git@github.com:sacloud/homebrew-packer-builder-sakuracloud.git homebrew-packer-builder-sakuracloud
cd homebrew-packer-builder-sakuracloud

cat << EOL > packer-builder-sakuracloud.rb
class PackerBuilderSakuracloud < Formula

  _version = "${VERSION}"
  sha256_src = "${SHA256_SRC}"

  desc "Packer builder plugin for SakuraCloud"
  homepage "https://github.com/sacloud/packer-builder-sakuracloud"
  url "https://github.com/sacloud/packer-builder-sakuracloud/releases/download/v#{_version}/packer-builder-sakuracloud_darwin-amd64.zip"
  sha256 sha256_src
  head "https://github.com/sacloud/packer-builder-sakuracloud.git"
  version _version

  depends_on "packer" => :run

  def install
    bin.install "packer-builder-sakuracloud"
  end

  def caveats; <<-EOS.undent

    This plugin requires locate into "~/.packer.d/plugins" directory.
    To enable, execute following:

        mkdir -p ~/.packer.d/plugins
        ln -s #{bin/"packer-builder-sakuracloud"} ~/.packer.d/plugins/

  EOS
  end


  test do
    minimal = testpath/"minimal.tf"
    minimal.write <<-EOS.undent
    {
      "builders": [{
          "type": "sakuracloud",
          "zone": "is1b",
          "os_type": "centos",
          "password": "this_is_fake_password"
      }]
    }
    EOS
    system "packer", "validate", minimal
  end
end
EOL

# show diff
git diff