module github.com/ddelnano/terraform-provider-mikrotik

go 1.16

require (
	github.com/ddelnano/terraform-provider-mikrotik/client v0.0.0-00010101000000-000000000000
	github.com/hashicorp/terraform-plugin-sdk/v2 v2.20.0
)

replace github.com/ddelnano/terraform-provider-mikrotik/client => ./client
