module github.com/ddelnano/terraform-provider-mikrotik

go 1.16

require (
	github.com/ddelnano/terraform-provider-mikrotik/client v0.0.0-00010101000000-000000000000
	github.com/golang/snappy v0.0.1 // indirect
	github.com/hashicorp/terraform-plugin-sdk/v2 v2.9.0
	github.com/pierrec/lz4 v2.0.5+incompatible // indirect
)

replace github.com/ddelnano/terraform-provider-mikrotik/client => ./client
