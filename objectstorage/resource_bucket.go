package objectstorage

import (
	"log"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/josajunior81/terraform-provider-vcd-object-storage-ext/pkg"
)

func resourceBucket() *schema.Resource {
	return &schema.Resource{
		Create: resourceBucketCreate,
		Read:   resourceBucketRead,
		Update: resourceBucketUpdate,
		Delete: resourceBucketDelete,

		Schema: map[string]*schema.Schema{
			"last_updated": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"tag": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: false,
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
		},
	}
}

// Creates Bucket on the Object Storage
func resourceBucketCreate(d *schema.ResourceData, meta interface{}) error {
	s3client := meta.(pkg.S3Client)

	d.SetId(uuid.NewString())
	bucketName := d.Get("name").(string)
	tags := d.Get("tag").([]any)

	s3client.CreateBucket(bucketName)
	s3client.BucketTags(bucketName, tags)

	return resourceBucketRead(d, meta)
}

// Reads Bucket from Object Storage
func resourceBucketRead(d *schema.ResourceData, meta interface{}) error {
	resourceID := d.Id()
	log.Println(">>> resourceID:", resourceID)
	return nil
}

func resourceBucketUpdate(d *schema.ResourceData, meta interface{}) error {
	// resourceID := d.Id()
	s3client := meta.(pkg.S3Client)

	bucketName := d.Get("name").(string)
	tags := d.Get("tag").([]any)

	s3client.BucketTags(bucketName, tags)
	return nil
}

// Deletes Bucket at the Object Storage
func resourceBucketDelete(d *schema.ResourceData, meta interface{}) error {
	resourceID := d.Id()
	log.Println(">>> resourceID:", resourceID)
	return nil
}
