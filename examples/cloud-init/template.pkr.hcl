packer {
  required_plugins {
    sakuracloud = {
      version = ">= 0.11.0"
      source = "github.com/sacloud/sakuracloud"
    }
  }
}

variable "ssh_public_key_path" {
  type    = string
  default = "~/.ssh/id_ed25519.pub"
}
locals {
  ssh_public_key = trimspace(file(pathexpand(var.ssh_public_key_path)))
}

source "sakuracloud" "example" {
  zone  = "is1b"

  # cloud-init対応のアーカイブの場合os_typeに"custom"を指定し、source_archiveにアーカイブIDを指定する
  os_type   = "custom"
  source_archive = 113601947034 # Ubuntu Server 24.04.1 LTS 64bit (cloudimg)

  disk_size = 20
  disk_plan = "ssd"

  user_data = <<-CLOUDINIT
    #cloud-config
    package_update: true
    packages:
      - curl
      - jq
    ssh_authorized_keys:
      - ${local.ssh_public_key}
    runcmd:
      - [ bash, -lc, "echo hello from cloud-init" ]
  CLOUDINIT

  # cloud-initを利用する場合以下の項目を指定する必要があります
  user_name = "ubuntu"
  disable_generate_public_key = true
  ssh_private_key_file = "~/.ssh/id_ed25519"

  core        = 2
  memory_size = 4

  archive_name        = "packer-example-cloud-init"
  archive_description = "description of archive"
}

build {
  sources = [
    "source.sakuracloud.example"
  ]
  provisioner "shell" {
    execute_command = "sudo -E bash -euxo pipefail '{{ .Path }}'"

    inline = [
      "echo 'hello!'",

      # cloud-init の実行済み状態を消す
      "if command -v cloud-init >/dev/null 2>&1; then cloud-init status --wait || true; cloud-init clean --logs; fi",

      # machine-id
      "truncate -s 0 /etc/machine-id",
      "rm -f /var/lib/dbus/machine-id || true",

      # SSH ホスト鍵
      "rm -f /etc/ssh/ssh_host_* || true",

      # DHCP リース等、引き継ぎたくないネットワーク系
      "rm -f /var/lib/dhcp/* || true",
      "rm -f /var/lib/NetworkManager/*lease* || true",
      "rm -f /etc/udev/rules.d/70-persistent-net.rules || true",

      # ホスト名
      "rm -f /etc/hostname || true",

      # 履歴・ログ・一時ファイル
      "rm -f /root/.bash_history || true",
      "rm -f /home/*/.bash_history || true",
      "journalctl --rotate || true",
      "journalctl --vacuum-time=1s || true",
      "rm -rf /tmp/* /var/tmp/* || true",

      "sync",
    ]
  }
}

