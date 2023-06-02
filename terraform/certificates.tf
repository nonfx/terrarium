resource "aws_acm_certificate" "main" {
  count = var.certificate_enabled ? 1 : 0
  domain_name       = "*.${var.domain_name}"
  validation_method = "DNS"
}

data "aws_route53_zone" "main" {
  count = var.certificate_enabled ? 1 : 0
  name         = var.domain_name
  private_zone = false
}

resource "aws_route53_record" "main" {
  for_each = {
    for dvo in flatten([
      for cert in aws_acm_certificate.main: cert.domain_validation_options
    ]): dvo.domain_name => {
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
  zone_id         = data.aws_route53_zone.main[0].zone_id
}

resource "aws_acm_certificate_validation" "main" {
  count = var.certificate_enabled ? 1 : 0
  certificate_arn         = aws_acm_certificate.main[count.index].arn
  validation_record_fqdns = [for record in aws_route53_record.main : record.fqdn]
}