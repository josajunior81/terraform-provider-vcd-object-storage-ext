---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "vcd-object-storage-ext Provider"
subcategory: ""
description: |-
  
---

# vcd-object-storage-ext Provider





<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `api_token` (String) The Api Token to access VCD
- `org` (String) The org (tenat) Object Storage
- `s3_url` (String) The S3 url for Object Storage
- `vcd_url` (String) The VCD url

### Optional

- `insecure` (Boolean) If set, VCDClient will permit unverifiable SSL certificates.
- `region` (String) The S3 region for Object Storage
