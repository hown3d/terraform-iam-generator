resource "aws_s3_bucket" "this" {
  bucket = "${var.hello}-${var.foo}-bucket"
}

variable "hello" {
}

variable "foo" {
}
