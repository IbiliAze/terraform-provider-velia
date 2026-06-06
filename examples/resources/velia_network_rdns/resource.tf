data "velia_network" "main" {
  filter_cidr = "192.168.190.96/28"
}

resource "velia_network_rdns" "example" {
  network_id = data.velia_network.main.id
  ip         = "192.168.190.100"
  type       = "PTR"
  rdata      = "www.example.com"
}
