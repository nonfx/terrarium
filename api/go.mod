module github.com/cldcvr/terrarium/api

go 1.20

replace github.com/cldcvr/terrarium/api/pkg/terraform-config-inspect => ./pkg/terraform-config-inspect

require (
	github.com/cldcvr/terrarium/api/pkg/terraform-config-inspect v0.0.0-00010101000000-000000000000
	github.com/rotisserie/eris v0.5.4
	gorm.io/driver/postgres v1.5.0
	gorm.io/gorm v1.25.1
)

require (
	github.com/agext/levenshtein v1.2.2 // indirect
	github.com/apparentlymart/go-textseg/v13 v13.0.0 // indirect
	github.com/google/go-cmp v0.3.1 // indirect
	github.com/hashicorp/hcl v0.0.0-20170504190234-a4b07c25de5f // indirect
	github.com/hashicorp/hcl/v2 v2.16.2 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgx/v5 v5.3.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/mitchellh/go-wordwrap v1.0.0 // indirect
	github.com/zclconf/go-cty v1.12.1 // indirect
	golang.org/x/crypto v0.6.0 // indirect
	golang.org/x/text v0.8.0 // indirect
)
