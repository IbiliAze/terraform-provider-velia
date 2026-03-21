#!/bin/bash

cd ..
go build -o /Users/ibi/Documents/git/go/providers/velia/terraform-provider-velia
cd test/

echo 'provider_installation {
  dev_overrides {
    "registry.terraform.io/ibiliaze/velia" = "/Users/ibi/Documents/git/go/providers/velia"
  }

  direct {}
}' > dev.tfrc

TF_CLI_CONFIG_FILE="$(pwd)/dev.tfrc" terraform init
TF_CLI_CONFIG_FILE="$(pwd)/dev.tfrc" terraform plan