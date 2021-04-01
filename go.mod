module github.com/ddelnano/terraform-provider-mikrotik

go 1.14

require (
	cloud.google.com/go/storage v1.0.0 // indirect
	github.com/ddelnano/terraform-provider-mikrotik/client v0.0.0-00010101000000-000000000000
	github.com/go-routeros/routeros v0.0.0-20210123142807-2a44d57c6730
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/hashicorp/terraform-plugin-sdk v1.15.0
	github.com/stretchr/testify v1.4.0 // indirect
	golang.org/x/tools v0.0.0-20200501155019-2658dc0cadb5 // indirect
	gopkg.in/yaml.v2 v2.2.8 // indirect
)

replace github.com/ddelnano/terraform-provider-mikrotik/client => ./client
