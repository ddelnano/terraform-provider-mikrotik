resource "mikrotik_dns_record" "record" {
  name    = "example.domain.com"
  address = "192.168.88.1"
  ttl     = 300
}

resource "mikrotik_dns_record" "record_regexp" {
  regexp  = ".+\\.example\\.domain\\.com"
  address = "192.168.88.1"
  ttl     = 300
}
