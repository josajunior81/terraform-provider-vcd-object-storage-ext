package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/josajunior81/terraform-provider-vcd-object-storage-ext/objectstorage"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: objectstorage.Provider,
	})
}
