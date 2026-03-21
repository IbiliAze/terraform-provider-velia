terraform {
  required_providers {
    velia = {
      source  = "eightmile/velia"
      version = "0.0.1"
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

output "contact_id" {
  value = velia_customer_contact.test.id
}