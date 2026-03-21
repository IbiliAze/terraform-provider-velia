package main

import (
	"context"
	"log"

	"github.com/IbiliAze/terraform-provider-velia/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

func main() {

	err := providerserver.Serve(context.Background(), provider.New, providerserver.ServeOpts{
		Address: "registry.terraform.io/eightmile/velia",
	})
	if err != nil {
		log.Fatal(err)
	}
}
