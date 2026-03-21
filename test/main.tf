terraform {
  required_providers {
    velia = {
      source  = "ibiliaze/velia"
      version = "0.1.2"
    }
  }
}

provider "velia" {
  endpoint  = "https://www.velia.net/api/v1"
  api_token = "YOUR_API_TOKEN"
}

resource "velia_customer_contact" "test" {
  email = "user@example.com"
  type  = "billing"
}

resource "velia_server_group" "test" {
  name = "23"
  color = "blue"
  servers = ["123"]
}

output "contact_id" {
  value = velia_customer_contact.test.id
}