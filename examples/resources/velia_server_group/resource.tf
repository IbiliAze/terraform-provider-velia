resource "velia_server_group" "example" {
  name    = "web-servers"
  color   = "#13355b"
  servers = [12345, 12346]
}