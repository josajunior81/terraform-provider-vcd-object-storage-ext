package pkg

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/vmware/go-vcloud-director/v2/govcd"
)

const PATH = "api/v1/s3"

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

func (s S3Client) doRequest(method, reqUrl, body string, additionalHeaders map[string]string) (string, error) {
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

	req.Header.Add("Authorization", "Bearer "+s.bearerToken)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("accept", "application/json")
	for k, v := range additionalHeaders {
		req.Header.Add(k, v)
	}

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
		log.Printf("Body: %s", string(body))
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
	bucketUrl := s.mountUrl(name, "max-keys=1")

	resp, err := s.doRequest(http.MethodGet, bucketUrl, "", nil)

	return resp, err
}

func (s S3Client) CreateBucket(name string) error {
	createBucketUrl := s.mountUrl(name, "")

	body := fmt.Sprintf(`{"name":"%s", "locationConstraint":"%s"}`, name, s.region)

	_, err := s.doRequest(http.MethodPut, createBucketUrl, body, nil)

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

	_, err := s.doRequest(http.MethodPut, tagsUrl, tagSet, nil)

	return err
}

func (s S3Client) removeBucketTags(bucket string) error {
	tagsUrl := s.mountUrl(bucket, "tagging")

	_, err := s.doRequest(http.MethodDelete, tagsUrl, "", nil)
	return err
}

func (s S3Client) UploadObject(bucket, key, source string, overwrite bool) error {
	objectUrl := s.mountUrl(bucket+"/"+key, fmt.Sprintf("overwrite=%t", overwrite))
	return s.doUpload(objectUrl, source)
}

func (s S3Client) BucketAcls(bucket string, setDefault bool, cannedAcl string, aclsI []interface{}) error {
	aclsUrl := s.mountUrl(bucket, "acl")

	bucketObj, err := s.getBucketObject(bucket)
	if err != nil {
		log.Panicf("ERROR getting bucker %v", err)
		return err
	}

	cannedAclHeader := make(map[string]string)
	if len(cannedAcl) > 0 && cannedAcl != "" {
		cannedAclHeader["X-Amz-Acl"] = cannedAcl
	}

	if setDefault {
		return s.defaultAcl(bucket, bucketObj, cannedAclHeader)
	}

	var grants []map[string]interface{}

	log.Printf("ACL %v", aclsI...)
	for _, a := range aclsI {
		acl := a.(map[string]interface{})

		log.Printf("ACL USER: %s", acl["user"])

		grantee := map[string]interface{}{}
		grant := map[string]interface{}{}
		switch acl["user"] {
		case "TENANT":
			grantee["id"] = bucketObj.Tenant + "|"
		case "AUTHENTICATED":
			grantee["uri"] = "http://acs.amazonaws.com/groups/global/AuthenticatedUsers"
		case "PUBLIC":
			grantee["uri"] = "http://acs.amazonaws.com/groups/global/AllUsers"
		case "SYSTEM-LOGGER":
			grantee["uri"] = "http://acs.amazonaws.com/groups/s3/LogDelivery"
		}
		grant["grantee"] = grantee
		grant["permission"] = acl["permission"]

		grants = append(grants, grant)
	}

	grants = append(grants, map[string]interface{}{"grantee": map[string]interface{}{"id": bucketObj.Owner.Id}, "permission": "FULL_CONTROL"})

	log.Printf("GRANTS %v", grants)

	payload := map[string]interface{}{}
	payload["owner"] = bucketObj.Owner
	payload["grants"] = grants

	log.Printf("payload %v", payload)

	payloadStr, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	log.Println("payload STR ==> " + string(payloadStr))

	_, err1 := s.doRequest(http.MethodPut, aclsUrl, string(payloadStr), cannedAclHeader)

	return err1
}

func toCamelCase(s string) string {
	re, _ := regexp.Compile(`[-_]\w`)
	res := re.ReplaceAllStringFunc(s, func(m string) string {
		return strings.ToUpper(m[1:])
	})
	return res
}

func (s S3Client) BucketCors(bucket string, corsI []interface{}) error {
	corsUrl := s.mountUrl(bucket, "cors")

	cors := map[string]interface{}{}

	payload := map[string][]interface{}{}

	for _, c := range corsI {
		for k, v := range c.(map[string]interface{}) {
			cors[toCamelCase(k)] = v
		}
		payload["corsRules"] = append(payload["corsRules"], cors)
	}

	// cors["id"] = uuid.New()

	log.Printf("payload %v", payload)

	payloadStr, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	log.Println("payload STR ==> " + string(payloadStr))

	_, err1 := s.doRequest(http.MethodPut, corsUrl, string(payloadStr), nil)

	return err1
}

func (s S3Client) defaultAcl(bucketName string, bucket *Bucket, cannedAclHeader map[string]string) error {
	aclsUrl := s.mountUrl(bucketName, "acl")

	var grants []map[string]interface{}

	grants = append(grants, map[string]interface{}{"grantee": map[string]interface{}{"id": bucket.Owner.Id}, "permission": "FULL_CONTROL"})

	payload := map[string]interface{}{}
	payload["owner"] = bucket.Owner
	payload["grants"] = grants

	payloadStr, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	log.Println(string(payloadStr))

	_, err1 := s.doRequest(http.MethodPut, aclsUrl, string(payloadStr), cannedAclHeader)

	return err1
}

func (s S3Client) getBucketObject(name string) (*Bucket, error) {
	bucketStr, err := s.GetBucket(name)
	if err != nil {
		log.Panicf("ERROR getting bucker %v", err)
		return nil, err
	}

	var bucketObj *Bucket
	if err := json.Unmarshal([]byte(bucketStr), &bucketObj); err != nil {
		log.Println(bucketStr)
		log.Panicf("ERROR unmarshalling bucket %v", err)
		return nil, err
	}

	return bucketObj, nil
}

func (s S3Client) DeleteBucket(name string) error {
	deleteUrl := s.mountUrl(name, "")

	payload := `{
		"quiet": true,
		"removeAll": true,
		"deleteVersion": true,
		"tryAsync": true
	}`

	_, err := s.doRequest(http.MethodPost, s.mountUrl(name, "delete"), payload, nil)
	if err != nil {
		log.Panicf("ERROR deleting all Objetcts of bucket %s: %v", name, err)
		return err
	}

	_, err2 := s.doRequest(http.MethodDelete, deleteUrl, "", nil)
	return err2
}
