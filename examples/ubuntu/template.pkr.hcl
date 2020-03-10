locals {
  password = "TestUserPassword01"
}


source "sakuracloud" "example" {
  zone = "is1b"

  os_type   = "ubuntu"
  user_name = "ubuntu"
  password  = local.password

  core        = 2
  memory_size = 4

  archive_name        = "packer-example-ubuntu"
  archive_description = "description of archive"
}

build {
  sources = [
    "source.sakuracloud.example"
  ]
  provisioner "shell" {
    execute_command = "echo '${local.password}' | {{ .Vars }} sudo -E -S sh '{{ .Path }}'"
    inline = [
      "apt-get update -y",
      "apt-get install -y curl"
    ]
  }
}

