resource "aws_s3_bucket" "dcalvo-dev-bucket" {
  bucket = var.bucket-name
  acl    = "private"

  website {
    index_document = "index.html"
    error_document = "error.html"
  }

  provisioner "local-exec" {
    when    = create
    command = "echo Hello from ${var.domain-name} > index.html && aws s3 cp index.html s3://${var.bucket-name} && rm index.html"
  }

  provisioner "local-exec" {
    when    = destroy
    command = "aws s3 rm s3://${self.bucket} --recursive"
  }
}

resource "aws_s3_bucket_policy" "dcalvo-dev-bucket-bucket-policy" {
  bucket = aws_s3_bucket.dcalvo-dev-bucket.id

  policy = <<POLICY
{
    "Version": "2008-10-17",
    "Id": "PolicyForCloudFrontPrivateContent",
    "Statement": [
        {
            "Sid": "1",
            "Effect": "Allow",
            "Principal": {
                "AWS": "arn:aws:iam::cloudfront:user/CloudFront Origin Access Identity ${aws_cloudfront_origin_access_identity.dcalvo-dev-origin-access-identity.id}"
            },
            "Action": "s3:GetObject",
            "Resource": "arn:aws:s3:::${aws_s3_bucket.dcalvo-dev-bucket.id}/*"
        }
    ]
}
POLICY
}

resource "aws_acm_certificate" "dcalvo-dev-cert" {
  //apparently if you change providers it will not delete the old resource, interesting
  provider          = aws.us-east-1
  domain_name       = var.domain-name
  validation_method = "DNS"

  tags = {
    Environment = var.domain-name
  }
  lifecycle {
    create_before_destroy = true
  }
}

//This works but it I have no knowledge on the details
resource "aws_route53_record" "acm-validation-route53-dcalvo-dev" {
  provider = aws.us-east-1
  for_each = {
    for dvo in aws_acm_certificate.dcalvo-dev-cert.domain_validation_options : dvo.domain_name => {
      name   = dvo.resource_record_name
      record = dvo.resource_record_value
      type   = dvo.resource_record_type
    }
  }

  allow_overwrite = true
  name            = each.value.name
  records         = [each.value.record]
  ttl             = 60
  type            = each.value.type
  zone_id         = data.terraform_remote_state.dcalvo-dev-zone.outputs.zone_id
}

resource "aws_acm_certificate_validation" "acm-validation-dcalvo-dev" {
  provider                = aws.us-east-1
  certificate_arn         = aws_acm_certificate.dcalvo-dev-cert.arn
  validation_record_fqdns = [for record in aws_route53_record.acm-validation-route53-dcalvo-dev : record.fqdn]
}

resource "aws_cloudfront_origin_access_identity" "dcalvo-dev-origin-access-identity" {
  comment = "Access identity for ${var.domain-name}"
}

resource "aws_cloudfront_distribution" "dcalvo-dev-distribution" {
  aliases = [var.domain-name]

  origin {
    domain_name = aws_s3_bucket.dcalvo-dev-bucket.bucket_regional_domain_name
    origin_id   = aws_s3_bucket.dcalvo-dev-bucket.id

    s3_origin_config {
      origin_access_identity = aws_cloudfront_origin_access_identity.dcalvo-dev-origin-access-identity.cloudfront_access_identity_path
    }
  }

  viewer_certificate {
    acm_certificate_arn = aws_acm_certificate.dcalvo-dev-cert.arn
    ssl_support_method  = "sni-only"
  }

  enabled             = true
  is_ipv6_enabled     = true
  comment             = "Cloudfront for ${var.domain-name} with S3"
  default_root_object = "index.html"

  //  logging_config {
  //    include_cookies = false
  //    bucket          = "mylogs.s3.amazonaws.com"
  //    prefix          = "myprefix"
  //  }

  default_cache_behavior {
    allowed_methods  = ["GET", "HEAD"]
    cached_methods   = ["GET", "HEAD"]
    target_origin_id = aws_s3_bucket.dcalvo-dev-bucket.id

    forwarded_values {
      query_string = false

      cookies {
        forward = "none"
      }
    }

    viewer_protocol_policy = "allow-all"
    min_ttl                = 0
    default_ttl            = 3600
    max_ttl                = 86400
  }

  restrictions {
    geo_restriction {
      restriction_type = "none"
    }
  }

  tags = {
    Environment = var.domain-name
  }
}


resource "aws_route53_record" "dcalvo-dev-a" {
  zone_id = data.terraform_remote_state.dcalvo-dev-zone.outputs.zone_id
  name    = var.domain-name
  type    = "A"

  alias {
    name                   = aws_cloudfront_distribution.dcalvo-dev-distribution.domain_name
    zone_id                = aws_cloudfront_distribution.dcalvo-dev-distribution.hosted_zone_id
    evaluate_target_health = true
  }
}

resource "null_resource" "populate-s3-bucket" {
  depends_on = [aws_s3_bucket.dcalvo-dev-bucket]
  provisioner "local-exec" {
    command = "aws s3 cp index.html s3://${var.bucket-name}"
  }
}