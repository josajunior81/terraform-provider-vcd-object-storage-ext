terraform {
  required_providers {
    vcd-object-storage-ext = {
      source = "josajunior81/vcd-object-storage-ext"
      version = "0.0.1"
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
}

resource "vcd-object-storage-ext_bucket" "this" {
  name = "provider-teste"
  dynamic "tag" {
    for_each = local.tags
    content {
      name = tag.value.name
      value = tag.value.value
    }
  }
}

resource "vcd-object-storage-ext_object" "this" {
  bucket = vcd-object-storage-ext_bucket.this.name
  key = "job.log"
  source = "/home/s827289255/Downloads/job.log"
}