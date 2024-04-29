package objectstorage

import "github.com/josajunior81/terraform-provider-vcd-object-storage-ext/pkg"

type Config struct {
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
	Url       string `json:"s3_url"`
	S3Client  pkg.S3Client
}
