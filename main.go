package main

//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

import (
	"context"
	"log"

	"github.com/IbiliAze/terraform-provider-velia/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

func main() {

	err := providerserver.Serve(context.Background(), provider.New, providerserver.ServeOpts{
		Address: "registry.terraform.io/ibiliaze/velia",
	})
	if err != nil {
		log.Fatal(err)
	}
}
