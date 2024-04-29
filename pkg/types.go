package pkg

import "net/http"

type S3Client struct {
	client      *http.Client
	s3Url       string
	region      string
	bearerToken string
	path        string
}

type Bucket struct {
	Name      string `json:"name"`
	Tenant    string `json:"tenant"`
	S3Href    string `json:"s3Href"`
	S3AltHref string `json:"s3AltHref"`
	Owner     Owner  `json:"owner"`
}

type Owner struct {
	Id          string `json:"id"`
	DisplayName string `json:"displayName"`
}
