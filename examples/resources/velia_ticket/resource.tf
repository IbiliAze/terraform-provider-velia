resource "velia_ticket" "example" {
  topic   = "velianet-support"
  subject = "Network connectivity issue on server 12345"
  message = "Since 10:00 UTC we are experiencing packet loss on eth0."
  servers = [12345]
}
