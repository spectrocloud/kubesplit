module github.com/spectrocloud/kubesplit

go 1.18

replace github.com/kairos-io/kairos => /home/dimitris/workspace/kairos/kairos

require (
	github.com/kairos-io/kairos v1.3.0
	github.com/onsi/ginkgo/v2 v2.5.1
	github.com/onsi/gomega v1.24.1
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/google/go-cmp v0.5.9 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/itchyny/gojq v0.12.10 // indirect
	github.com/itchyny/timefmt-go v0.1.5 // indirect
	github.com/kr/pretty v0.1.0 // indirect
	golang.org/x/net v0.2.0 // indirect
	golang.org/x/sys v0.2.0 // indirect
	golang.org/x/text v0.4.0 // indirect
	gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127 // indirect
)

replace gopkg.in/yaml.v3 => github.com/atomatt/yaml v0.0.0-20200228174225-55c5cf55e3ee
