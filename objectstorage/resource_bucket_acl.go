package objectstorage

// import (
// 	"log"

// 	"github.com/google/uuid"
// 	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
// 	"github.com/josajunior81/terraform-provider-vcd-object-storage-ext/pkg"
// )

// func resourceBucketAcl() *schema.Resource {
// 	return &schema.Resource{
// 		Create: resourceBucketAclCreate,
// 		Read:   resourceBucketAclRead,
// 		Update: resourceBucketAclUpdate,
// 		Delete: resourceBucketAclDelete,

// 		Schema: map[string]*schema.Schema{
// 			"last_updated": {
// 				Type:     schema.TypeString,
// 				Computed: true,
// 			},
// 			"bucket": {
// 				Type:     schema.TypeString,
// 				Required: true,
// 				ForceNew: true,
// 			},

// 			"acl": {
// 				Type:     schema.TypeList,
// 				Optional: true,
// 				ForceNew: false,
// 				Elem: &schema.Resource{
// 					Schema: map[string]*schema.Schema{
// 						"user": {
// 							Type:     schema.TypeString,
// 							Required: true,
// 						},
// 						"permission": {
// 							Type:     schema.TypeString,
// 							Required: true,
// 						},
// 					},
// 				},
// 			},
// 		},
// 	}
// }

// // Creates Bucket on the Object Storage
// func resourceBucketAclCreate(d *schema.ResourceData, meta interface{}) error {
// 	s3client := meta.(pkg.S3Client)

// 	d.SetId(uuid.NewString())
// 	bucketName := d.Get("name").(string)
// 	tags := d.Get("tag").([]any)

// 	s3client.CreateBucket(bucketName)
// 	s3client.BucketTags(bucketName, tags)

// 	return resourceBucketAclRead(d, meta)
// }

// // Reads Bucket from Object Storage
// func resourceBucketAclRead(d *schema.ResourceData, meta interface{}) error {
// 	resourceID := d.Id()
// 	log.Println(">>> resourceID:", resourceID)
// 	return nil
// }

// func resourceBucketAclUpdate(d *schema.ResourceData, meta interface{}) error {
// 	// resourceID := d.Id()
// 	s3client := meta.(pkg.S3Client)

// 	bucketName := d.Get("name").(string)
// 	tags := d.Get("tag").([]any)

// 	s3client.BucketTags(bucketName, tags)
// 	return nil
// }

// // Deletes Bucket at the Object Storage
// func resourceBucketAclDelete(d *schema.ResourceData, meta interface{}) error {
// 	resourceID := d.Id()
// 	log.Println(">>> resourceID:", resourceID)
// 	return nil
// }
