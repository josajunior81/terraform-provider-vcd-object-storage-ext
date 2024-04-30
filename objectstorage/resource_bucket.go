package objectstorage

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/josajunior81/terraform-provider-vcd-object-storage-ext/pkg"
)

func resourceBucket() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBucketCreate,
		ReadContext:   resourceBucketRead,
		UpdateContext: resourceBucketUpdate,
		DeleteContext: resourceBucketDelete,
		Description:   "A bucket is a container for storing objects in a compartment within an Object Storage namespace.",

		Schema: map[string]*schema.Schema{
			"last_updated": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The bucket name. It must be URL encoded.",
			},

			"canned_acl": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
				ValidateDiagFunc: func(v interface{}, p cty.Path) diag.Diagnostics {
					value := v.(string)
					var diags diag.Diagnostics

					switch value {
					case "private", "public-read", "public-read-write", "authenticated-read", "group-read-write", "group-read", "log-delivery-write":
						return diags
					default:
						diag := diag.Diagnostic{
							Severity:      diag.Error,
							Summary:       "Wrong value. Valid Values: private | public-read | public-read-write | authenticated-read | group-read-write | group-read | log-delivery-write",
							Detail:        fmt.Sprintf("%q is not x-amz-acl valid value", value),
							AttributePath: p,
						}

						return append(diags, diag)
					}
				},
				Description: "The ACL of the bucket using the specified canned ACL. Valid Values: private | public-read | public-read-write | authenticated-read.",
			},

			"tag": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    false,
				Description: "The bucket tags.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},

			"acl": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    false,
				Description: "Access control lists (ACLs) enable you to manage access to buckets and objects.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user": {
							Type:     schema.TypeString,
							Required: true,
							ValidateDiagFunc: func(v interface{}, p cty.Path) diag.Diagnostics {
								value := v.(string)
								var diags diag.Diagnostics

								if value != "TENANT" && value != "AUTHENTICATED" && value != "PUBLIC" && value != "SYSTEM-LOGGER" {
									diag := diag.Diagnostic{
										Severity:      diag.Error,
										Summary:       "Wrong value. Valid Values: TENANT | AUTHENTICATED | PUBLIC | SYSTEM-LOGGER",
										Detail:        fmt.Sprintf("%q is not a valid ACL User", value),
										AttributePath: p,
									}

									diags = append(diags, diag)
								}

								return diags
							},
							Description: "ACL users. Valid Values: TENANT | AUTHENTICATED | PUBLIC | SYSTEM-LOGGER",
						},
						"permission": {
							Type:     schema.TypeString,
							Required: true,
							ValidateDiagFunc: func(i interface{}, p cty.Path) diag.Diagnostics {
								value := i.(string)
								var diags diag.Diagnostics

								if value != "FULL_CONTROL" && value != "READ" && value != "WRITE" && value != "READ_ACP" && value != "WRITE_ACP" {
									diag := diag.Diagnostic{
										Severity:      diag.Error,
										Summary:       "Wrong value. Valid Values: FULL_CONTROL | READ | WRITE | READ_ACP | WRITE_ACP",
										Detail:        fmt.Sprintf("%q is not a valid ACL Permission", value),
										AttributePath: p,
									}

									diags = append(diags, diag)
								}

								return diags
							},
							Description: "ACL permission. Valid Values: FULL_CONTROL | READ | WRITE | READ_ACP | WRITE_ACP",
						},
					},
				},
			},

			"cors": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    false,
				Description: "Cross-origin resource sharing (CORS) defines a way for client web applications that are loaded in one domain to interact with resources in a different domain.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"allowed_headers": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "The allowed_headers element specifies which headers are allowed in a preflight request through the Access-Control-Request-Headers header. Must be a comma separated string",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"expose_headers": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Each expose_headers element identifies a header in the response that you want customers to be able to access from their applications (for example, from a JavaScript XMLHttpRequest object). Must be a comma separated string",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"allowed_methods": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "In the CORS configuration, you can specify the following values for the allowed_methods element GET | PUT | POST | DELETE | HEAD. Must be a comma separated string",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"allowed_origins": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "In the allowed_origins element, you specify the origins that you want to allow cross-domain requests from. Must be a comma separated string",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"max_age_seconds": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     3600,
							Description: "Max age in secods. Default 3600",
						},
					},
				},
			},
		},
	}
}

// Creates Bucket on the Object Storage
func resourceBucketCreate(c context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	s3client := meta.(pkg.S3Client)

	var diags diag.Diagnostics

	d.SetId(uuid.NewString())
	bucketName := d.Get("name").(string)

	err := s3client.CreateBucket(bucketName)
	if err != nil {
		diag := diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error creating bucket",
			Detail:   fmt.Sprintf("Error creating bucket: %v", err),
		}
		return append(diags, diag)
	}

	diags = resourceBucketUpdate(c, d, meta)

	resourceBucketRead(c, d, meta)
	return diags
}

// Reads Bucket from Object Storage
func resourceBucketRead(c context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	resourceID := d.Id()
	log.Println(">>> resourceID:", resourceID)
	return nil
}

func resourceBucketUpdate(c context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// resourceID := d.Id()
	s3client := meta.(pkg.S3Client)

	var diags diag.Diagnostics

	bucketName := d.Get("name").(string)
	cannedAcl := d.Get("canned_acl").(string)
	tags := d.Get("tag").([]any)
	acls := d.Get("acl").([]interface{})
	cors := d.Get("cors").([]interface{})

	log.Printf("cors %v", cors)

	if cannedAcl != "" && len(acls) > 0 {
		diag := diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "You have to choose between Canned ACL or ACL rules",
			Detail:   "Choose between Canned ACL or ACL rules",
		}
		return append(diags, diag)
	}

	if len(acls) > 0 {
		err := s3client.BucketAcls(bucketName, false, cannedAcl, acls)
		if err != nil {
			diag := diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Error editing bucket ACLs",
				Detail:   fmt.Sprintf("Error editing bucket ACLs: %v", err),
			}
			return append(diags, diag)
		}
	} else {
		err := s3client.BucketAcls(bucketName, true, cannedAcl, nil)
		if err != nil {
			diag := diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Error editing bucket ACLs",
				Detail:   fmt.Sprintf("Error editing bucket ACLs: %v", err),
			}
			return append(diags, diag)
		}
	}

	if len(tags) > 0 {
		err := s3client.BucketTags(bucketName, tags)
		if err != nil {
			diag := diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Error editing bucket TAGs",
				Detail:   fmt.Sprintf("Error editing bucket TAGs: %v", err),
			}
			return append(diags, diag)
		}
	}

	if len(cors) > 0 {
		err := s3client.BucketCors(bucketName, cors)
		if err != nil {
			diag := diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Error editing bucket CORs",
				Detail:   fmt.Sprintf("Error editing bucket CORs: %v", err),
			}
			return append(diags, diag)
		}
	}
	return diags
}

// Deletes Bucket at the Object Storage
func resourceBucketDelete(c context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	resourceID := d.Id()
	log.Println(">>> resourceID:", resourceID)

	var diags diag.Diagnostics

	s3client := meta.(pkg.S3Client)
	bucketName := d.Get("name").(string)
	if err := s3client.DeleteBucket(bucketName); err != nil {
		diag := diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error editing bucket TAGs",
			Detail:   fmt.Sprintf("Error editing bucket TAGs: %v", err),
		}
		return append(diags, diag)
	}
	return diags
}
