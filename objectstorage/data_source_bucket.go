package objectstorage

import (
	"encoding/json"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/josajunior81/terraform-provider-vcd-object-storage-ext/pkg"
)

func dataSourceBucket() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBucketRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func dataSourceBucketRead(d *schema.ResourceData, meta interface{}) error {
	s3 := meta.(pkg.S3Client)
	name := d.Get("name").(string)

	bucket, err := s3.GetBucket(name)
	if err != nil {
		log.Printf("Error reading Bucket: %s", err)
		return err
	}

	var jsonBucket pkg.Bucket

	if err := json.Unmarshal([]byte(bucket), &jsonBucket); err != nil {
		log.Printf("Error Unmarshal Bucket: %s", err)
		return err
	}
	log.Printf("jsonBucket %v", jsonBucket)

	if err := d.Set("name", jsonBucket.Name); err != nil {
		log.Printf("Error d.Set(name, jsonBucket.Name) Bucket: %s", err)
		return err
	}

	return nil
}
