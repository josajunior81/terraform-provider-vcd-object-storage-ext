package pkg

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/vmware/go-vcloud-director/v2/govcd"
)

const PATH = "api/v1/s3"

type S3Client struct {
	client      *http.Client
	s3Url       string
	region      string
	bearerToken string
	path        string
}

func NewS3Client(s3url, region, apiToken, org, vcdUrl string) S3Client {

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	u, _ := url.ParseRequestURI(fmt.Sprintf("%s/api", vcdUrl))

	vcdClient := govcd.NewVCDClient(*u, true)
	_, err := vcdClient.SetApiToken(org, apiToken)
	if err != nil {
		log.Fatal(err)
	}

	bearerToken, err := vcdClient.GetBearerTokenFromApiToken(org, apiToken)
	if err != nil {
		log.Fatal(err)
	}

	s3client := S3Client{
		client:      client,
		s3Url:       s3url,
		bearerToken: bearerToken.AccessToken,
		path:        PATH,
		region:      region,
	}

	return s3client
}

func (s S3Client) mountUrl(resource, query string) string {
	if query != "" {
		return "https://" + s.s3Url + "/" + s.path + "/" + resource + "?" + query
	}
	return "https://" + s.s3Url + "/" + s.path + "/" + resource
}

func (s S3Client) doRequest(method, reqUrl, body string) (string, error) {
	var req *http.Request
	var err error
	if body != "" {
		req, err = http.NewRequest(method, reqUrl, bytes.NewBuffer([]byte(body)))
	} else {
		req, err = http.NewRequest(method, reqUrl, nil)
	}
	if err != nil {
		log.Printf("Error do request %v", err)
		return "", err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", s.bearerToken))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("accept", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		log.Printf("Error sending HTTP request: %s", err)
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Error reading HTTP response body: %s", err)
			return "", err
		}

		return string(body), nil
	}

	return "", nil
}

func (s S3Client) doUpload(reqUrl, source string) error {

	file, err := os.ReadFile(source)
	if err != nil {
		log.Println(err)
		return err
	}

	contentType := http.DetectContentType(file)

	log.Println("File content type", contentType)

	req, err := http.NewRequest(http.MethodPut, reqUrl, bytes.NewBuffer(file))
	if err != nil {
		log.Printf("Error do request %v", err)
		return err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", s.bearerToken))
	req.Header.Add("Content-Type", contentType)

	resp, err := s.client.Do(req)
	if err != nil {
		log.Printf("Error sending HTTP request: %s", err)
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (s S3Client) GetBucket(name string) (string, error) {
	bucketUrl := s.mountUrl(name, "")

	resp, err := s.doRequest(http.MethodGet, bucketUrl, "")

	return resp, err
}

func (s S3Client) CreateBucket(name string) error {
	createBucketUrl := s.mountUrl(name, "")

	body := fmt.Sprintf(`{"name":"%s", "locationConstraint":"%s"}`, name, s.region)

	_, err := s.doRequest(http.MethodPut, createBucketUrl, body)

	return err
}

func (s S3Client) BucketTags(bucket string, tags []any) error {
	tagsUrl := s.mountUrl(bucket, "tagging")

	s.removeBucketTags(bucket)

	tagSet := `{"tagSets":[ {"tags":[`
	for i := 0; i < len(tags); i++ {
		obj := tags[i].(map[string]interface{})
		tagSet = tagSet + fmt.Sprintf(`{"key":"%s","value":"%s"}`, obj["name"], obj["value"])
		if i+1 < len(tags) {
			tagSet = tagSet + ","
		}
	}
	tagSet = tagSet + `]}]}`

	_, err := s.doRequest(http.MethodPut, tagsUrl, tagSet)

	return err

}

func (s S3Client) removeBucketTags(bucket string) error {
	tagsUrl := s.mountUrl(bucket, "tagging")

	_, err := s.doRequest(http.MethodDelete, tagsUrl, "")
	return err
}

func (s S3Client) UploadObject(bucket, key, source string, overwrite bool) error {
	objectUrl := s.mountUrl(bucket+"/"+key, fmt.Sprintf("overwrite=%t", overwrite))
	return s.doUpload(objectUrl, source)
}