echo "provider_installation {
  dev_overrides {
    "registry.terraform.io/eightmile/velia" = "$(pwd)"
  }

  direct {}
}" > .terraform.rc