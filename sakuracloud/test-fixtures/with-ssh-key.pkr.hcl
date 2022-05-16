locals {
  password = "TestUserPassword01"
}

source "sakuracloud" "example" {
  zone = "is1a"

  os_type   = "ubuntu2004"
  user_name = "ubuntu"
  password  = local.password

  core        = 2
  memory_size = 4

  archive_name        = "packer-acctest-sshkey"
  archive_description = "description of archive"

  ssh_private_key_file = "%s"
}

build {
  sources = [
    "source.sakuracloud.example"
  ]
}

