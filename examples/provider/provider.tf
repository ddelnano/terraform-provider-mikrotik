# Configure the mikrotik Provider
provider "mikrotik" {
  host           = "hostname-of-server:8728"     # Or set MIKROTIK_HOST environment variable
  username       = "<username>"                  # Or set MIKROTIK_USER environment variable
  password       = "<password>"                  # Or set MIKROTIK_PASSWORD environment variable
  tls            = true                          # Or set MIKROTIK_TLS environment variable
  ca_certificate = "/path/to/ca/certificate.pem" # Or set MIKROTIK_CA_CERTIFICATE environment variable
  insecure       = true                          # Or set MIKROTIK_INSECURE environment variable
}
