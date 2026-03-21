echo "provider_installation {
  dev_overrides {
    "registry.terraform.io/ibiliaze/velia" = "$(pwd)"
  }

  direct {}
}" > .terraform.rc