packer {
  required_plugins {
    sakuracloud = {
      version = ">= 0.7.0"
      source  = "github.com/sacloud/sakuracloud"
    }
  }
}

variable "ssh_password" {
  default = "TestUserPassword01"
}

source "sakuracloud" "alpine" {
  zone                = "is1b"

  os_type  = "iso"
  iso_url  = "https://dl-cdn.alpinelinux.org/alpine/v3.23/releases/x86_64/alpine-virt-3.23.3-x86_64.iso"
  iso_checksum = "file:https://dl-cdn.alpinelinux.org/alpine/v3.23/releases/x86_64/alpine-virt-3.23.3-x86_64.iso.sha256"

  disk_size    = 20
  disk_plan    = "ssd"
  core         = 2
  memory_size  = 4

  archive_name        = "packer-alpine-example"
  archive_description = "Alpine Linux 3.23.3 (virt) image built with Packer"
  archive_tags        = ["alpine", "iso", "example"]

  boot_wait = "30s"
  us_keyboard = true
  boot_command = [
    # root でログイン（パスワードなし）
    "root<enter><wait5s>",

    # setup 開始
    "setup-alpine<enter><wait2s>",
    "us<enter><wait2s>", # keyboard layout は US を選択
    "us<enter><wait2s>", # keyboard layout は US を選択
    "<enter><wait10s>", # hostname はデフォルトのまま Enter
    "<enter><wait2s>", # interface は eth0 のまま Enter
    "dhcp<enter><wait2s>", # ip address は dhcp 取得
    "<enter><wait2s>", # それ以外のインターフェースは不要
    "${var.ssh_password}<enter><wait2s>", # root パスワード入力
    "${var.ssh_password}<enter><wait2s>",
    "<enter><wait2s>", # タイムゾーンは UTC のままで OK
    "<enter><wait2s>", # Proxy URL は変更しなくて OK
    "<enter><wait2s>", # APK URL もそのままで OK
    "<enter><wait2s>", # ユーザー作成不要
    "<enter><wait2s>", # SSH Server は openssh で OK
    "yes<enter><wait2s>", # Root login を許容(初回起動時のみ)
    "<enter><wait2s>", # 鍵登録しない
    "vda<enter><wait2s>", # disk は vda
    "sys<enter><wait2s>", # system disk にする
    "y<enter><wait2s>", # Erase disk and continue
    "<enter><wait2s>", # apk cache directory もデフォルト

    # 再起動
    "<wait30s>",
    "reboot<enter>",
  ]

  ssh_username = "root"
  ssh_password = var.ssh_password
  ssh_timeout  = "10m"
}

build {
  sources = ["source.sakuracloud.alpine"]

  provisioner "shell" {
    inline = [
      # PermitRootLogin を消す。
      "sed -i '/^PermitRootLogin yes$/d' /etc/ssh/sshd_config",
      "echo 'Alpine installed successfully'",
      "uname -a",
      "apk add curl",
    ]
  }
}
