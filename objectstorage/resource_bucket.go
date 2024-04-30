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
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         false,
				ValidateDiagFunc: validateCannedAcl,
				Description:      "The ACL of the bucket using the specified canned ACL. Valid Values: private | public-read | public-read-write | authenticated-read.",
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
				Description: "Access control lists (ACLs) enable you to manage access to buckets and objects",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validateAclUser,
							Description:      "ACL users. Valid Values: TENANT | AUTHENTICATED | PUBLIC | SYSTEM-LOGGER",
						},
						"permission": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validateAclPermission,
							Description:      "ACL permission. Valid Values: FULL_CONTROL | READ | WRITE | READ_ACP | WRITE_ACP",
						},
					},
				},
			},
		},
	}
}

func validateAclPermission(v interface{}, path cty.Path) diag.Diagnostics {
	value := v.(string)
	var diags diag.Diagnostics

	if value != "FULL_CONTROL" && value != "READ" && value != "WRITE" && value != "READ_ACP" && value != "WRITE_ACP" {
		diag := diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       "Wrong value. Valid Values: FULL_CONTROL | READ | WRITE | READ_ACP | WRITE_ACP",
			Detail:        fmt.Sprintf("%q is not a valid ACL Permission", value),
			AttributePath: path,
		}

		diags = append(diags, diag)
	}

	return diags
}

func validateAclUser(v interface{}, path cty.Path) diag.Diagnostics {
	value := v.(string)
	var diags diag.Diagnostics

	if value != "TENANT" && value != "AUTHENTICATED" && value != "PUBLIC" && value != "SYSTEM-LOGGER" {
		diag := diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       "Wrong value. Valid Values: TENANT | AUTHENTICATED | PUBLIC | SYSTEM-LOGGER",
			Detail:        fmt.Sprintf("%q is not a valid ACL User", value),
			AttributePath: path,
		}

		diags = append(diags, diag)
	}

	return diags
}

func validateCannedAcl(v interface{}, path cty.Path) diag.Diagnostics {
	value := v.(string)
	var diags diag.Diagnostics

	if value != "private" && value != "public-read" && value != "public-read-write" && value != "authenticated-read" {
		diag := diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       "Wrong value. Valid Values: private | public-read | public-read-write | authenticated-read",
			Detail:        fmt.Sprintf("%q is not x-amz-acl valid value", value),
			AttributePath: path,
		}

		diags = append(diags, diag)
	}

	return diags
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
	return diags
}

// Deletes Bucket at the Object Storage
func resourceBucketDelete(c context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	resourceID := d.Id()
	log.Println(">>> resourceID:", resourceID)
	return diag.Diagnostics{}
}
