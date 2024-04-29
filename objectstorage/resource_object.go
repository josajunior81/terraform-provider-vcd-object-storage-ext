package objectstorage

import (
	"log"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/josajunior81/terraform-provider-vcd-object-storage-ext/pkg"
)

func resourceObject() *schema.Resource {
	return &schema.Resource{
		Create: resourceObjectCreate,
		Read:   resourceObjectRead,
		Update: resourceObjectUpdate,
		Delete: resourceObjectDelete,

		Schema: map[string]*schema.Schema{
			"last_updated": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"bucket": {
				Type:     schema.TypeString,
				Required: true,
			},

			"key": {
				Type:     schema.TypeString,
				Required: true,
			},

			"source": {
				Type:     schema.TypeString,
				Required: true,
			},

			"overwrite": {
				Type:     schema.TypeBool,
				Computed: false,
				Optional: true,
				Default:  true,
			},

			// "tag": {
			// 	Type:     schema.TypeList,
			// 	Optional: true,
			// 	ForceNew: false,
			// 	Elem: &schema.Resource{
			// 		Schema: map[string]*schema.Schema{
			// 			"name": {
			// 				Type:     schema.TypeString,
			// 				Required: true,
			// 			},
			// 			"value": {
			// 				Type:     schema.TypeString,
			// 				Required: true,
			// 			},
			// 		},
			// 	},
			// },
		},
	}
}

// Creates Bucket on the Object Storage
func resourceObjectCreate(d *schema.ResourceData, meta interface{}) error {
	s3client := meta.(pkg.S3Client)

	d.SetId(uuid.NewString())
	bucket := d.Get("bucket").(string)
	key := d.Get("key").(string)
	source := d.Get("source").(string)
	overwrite := d.Get("overwrite").(bool)

	s3client.UploadObject(bucket, key, source, overwrite)

	return resourceObjectRead(d, meta)
}

// Reads Bucket from Object Storage
func resourceObjectRead(d *schema.ResourceData, meta interface{}) error {
	resourceID := d.Id()
	log.Println(">>> resourceID:", resourceID)
	return nil
}

func resourceObjectUpdate(d *schema.ResourceData, meta interface{}) error {
	// resourceID := d.Id()
	s3client := meta.(pkg.S3Client)

	bucketName := d.Get("name").(string)
	tags := d.Get("tag").([]any)

	s3client.BucketTags(bucketName, tags)
	return nil
}

// Deletes Bucket at the Object Storage
func resourceObjectDelete(d *schema.ResourceData, meta interface{}) error {
	resourceID := d.Id()
	log.Println(">>> resourceID:", resourceID)
	return nil
}
