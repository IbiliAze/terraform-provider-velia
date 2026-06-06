# Terraform Provider for Velia

This provider allows you to manage resources on [velia.net](https://www.velia.net) using Terraform.

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.21 (to build from source)

## Usage

```terraform
terraform {
  required_providers {
    velia = {
      source  = "ibiliaze/velia"
      version = "~> 0.3"
    }
  }
}

provider "velia" {
  api_token = var.velia_api_token
}
```

An API token can be generated in the [Velia customer portal](https://clients.velia.net).

## Resources

| Resource | Description |
|---|---|
| `velia_customer_contact` | Contact email address on the customer account |
| `velia_server_group` | Group of servers with a name and colour |
| `velia_server_label` | Label on an existing server |
| `velia_network_rdns` | Reverse DNS entry for a network IP |
| `velia_ticket` | Support ticket (no API delete — destroy removes from state only) |

## Data Sources

| Data Source | Description |
|---|---|
| `velia_customer_contact` | Look up a contact by ID |
| `velia_server_group` | Look up a server group by ID |
| `velia_server` | Look up a server by ID |
| `velia_network` | Look up a network by CIDR or IP |

## Example

```terraform
provider "velia" {
  api_token = var.velia_api_token
}

data "velia_server" "web" {
  id = 12345
}

resource "velia_server_label" "web" {
  server_id = data.velia_server.web.id
  label     = "web-01"
}

data "velia_network" "main" {
  filter_ip = data.velia_server.web.server_ip[0]
}

resource "velia_network_rdns" "web" {
  network_id = data.velia_network.main.id
  ip         = data.velia_server.web.server_ip[0]
  type       = "PTR"
  rdata      = "web-01.example.com"
}

resource "velia_server_group" "web" {
  name    = "web-servers"
  color   = "#13355b"
  servers = [data.velia_server.web.id]
}
```

## Building from Source

```bash
git clone https://github.com/IbiliAze/terraform-provider-velia
cd terraform-provider-velia
go build ./...
```

## License

[Mozilla Public License 2.0](https://www.mozilla.org/en-US/MPL/2.0/)
