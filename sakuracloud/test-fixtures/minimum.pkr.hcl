locals {
  password = "TestUserPassword01"
}

source "sakuracloud" "example" {
  zone = "is1a"

  os_type   = "ubuntu2404"
  user_name = "ubuntu"
  password  = local.password

  core        = 2
  memory_size = 4

  archive_name        = "packer-acctest-minimum"
  archive_description = "description of archive"
}

build {
  sources = [
    "source.sakuracloud.example"
  ]
}

