package main

import (
	"context"
	"log"

	"github.com/chriskuchin/terraform-provider-hiera5/hiera5"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

// Run "go generate" to format example terraform files and generate the docs for the registry/website

// If you do not have terraform installed, you can remove the formatting command, but its suggested to
// ensure the documentation is formatted properly.
//go:generate terraform fmt -recursive ./examples/

// Run the docs generation tool, check its repository for more information on how it works and how docs
// can be customized.
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

func main() {
	err := providerserver.Serve(
		context.Background(),
		hiera5.New,
		providerserver.ServeOpts{
			Address: "registry.terraform.io/chriskuchin/hiera5",
			// Debug:   true,
		},
	)

	if err != nil {
		log.Fatal(err)
	}
}
