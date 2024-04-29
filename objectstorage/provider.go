package objectstorage

import (
	// Documentation:
	// https://pkg.go.dev/github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema
	//

	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/josajunior81/terraform-provider-vcd-object-storage-ext/pkg"
)

var globalResourceMap = map[string]*schema.Resource{
	"vcd-object-storage-ext_bucket": resourceBucket(),
	"vcd-object-storage-ext_object": resourceObject(),
}

var globalDataSourceMap = map[string]*schema.Resource{
	"vcd-object-storage-ext_bucket": dataSourceBucket(),
}

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"s3_url": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("S3_URL", nil),
				Description: "The S3 url for Object Storage",
			},

			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("S3_REGION", nil),
				Description: "The S3 region for Object Storage",
			},

			"org": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("ORG", nil),
				Description: "The org (tenat) Object Storage",
			},

			"api_token": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("API_TOKEN", nil),
				Description: "The Api Token to access VCD",
			},

			"vcd_url": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("API_TOKEN", nil),
				Description: "The VCD url",
			},

			"insecure": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("INSECURE", false),
				Description: "If set, VCDClient will permit unverifiable SSL certificates.",
			},
		},
		ResourcesMap:         globalResourceMap,
		DataSourcesMap:       globalDataSourceMap,
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(_ context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	var providerDiag = diag.Diagnostics{}

	s3client := pkg.NewS3Client(d.Get("s3_url").(string), d.Get("region").(string), d.Get("api_token").(string), d.Get("org").(string), d.Get("vcd_url").(string))

	return s3client, providerDiag
}
