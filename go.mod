module github.com/spectrocloud/kubesplit

go 1.18

require (
	github.com/onsi/ginkgo/v2 v2.1.6
	github.com/onsi/gomega v1.20.2
	github.com/vmware-labs/go-yaml-edit v0.3.0
	github.com/vmware-labs/yaml-jsonpointer v0.1.0
	golang.org/x/text v0.3.7
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/go-openapi/jsonpointer v0.19.3 // indirect
	github.com/go-openapi/swag v0.19.5 // indirect
	github.com/google/go-cmp v0.5.8 // indirect
	github.com/mailru/easyjson v0.0.0-20190626092158-b2ccc519800e // indirect
	golang.org/x/net v0.0.0-20220722155237-a158d28d115b // indirect
	golang.org/x/sys v0.0.0-20220722155257-8c9f86f7a55f // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

replace gopkg.in/yaml.v3 => github.com/atomatt/yaml v0.0.0-20200228174225-55c5cf55e3ee
