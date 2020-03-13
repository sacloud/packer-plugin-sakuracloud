source "sakuracloud" "example" {
  zone  = "is1b"
  zones = ["is1a", "is1b", "tk1a", "tk1v"]

  os_type   = "centos8"
  password  = "TestUserPassword01"
  disk_size = 20
  disk_plan = "ssd"

  core        = 2
  memory_size = 4

  archive_name        = "packer-example-centos"
  archive_description = "description of archive"
}

build {
  sources = [
    "source.sakuracloud.example"
  ]
  provisioner "shell" {
    inline = [
      "echo 'hello!'",
    ]
  }
}

