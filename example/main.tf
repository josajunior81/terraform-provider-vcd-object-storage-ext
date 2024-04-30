terraform {
  required_providers {
    vcd-object-storage-ext = {
      source = "josajunior81/vcd-object-storage-ext"
      version = "0.0.2-alpha"
    }
  }
}

provider "vcd-object-storage-ext" {
  s3_url = var.s3_url
  org = var.org
  region = var.region
  api_token = var.api_token
  vcd_url = var.vcd_url
}

locals {
  tags = [{name="tag1", value="abc"}]
  acls = [{user="TENANT", permission="FULL_CONTROL"}, {user="AUTHENTICATED", permission="READ"}]
  cors = [{allowedMethods="GET,POST", maxAgeSeconds=3600, exposeHeaders=[], }]
}

resource "vcd-object-storage-ext_bucket" "this" {
  name = "provider-teste-2"

  # canned_acl = "public-read-write"

  # dynamic "tag" {
  #   for_each = local.tags
  #   content {
  #     name = tag.value.name
  #     value = tag.value.value
  #   }
  # }

  # dynamic "acl" {
  #   for_each = local.acls
  #   content {
  #     user = acl.value.user
  #     permission = acl.value.permission
  #   }
  # }  

  cors {
    allowed_methods = ["GET","POST"]
    // allowed_headers = []
    allowed_origins = ["*"]
    // expose_headers = []
    max_age_seconds = 3600
  }
}

# resource "vcd-object-storage-ext_object" "this" {
#   bucket = vcd-object-storage-ext_bucket.this.name
#   key = "job.log"
#   source = "/home/s827289255/Downloads/job.log"
# }