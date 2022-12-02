package parser_test

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/spectrocloud/kubesplit/internal/helm"
	. "github.com/spectrocloud/kubesplit/internal/parser"
)

var _ = Describe("Parser", func() {
	Context("Splitting elements", func() {
		It("Splits them correctly by type", func() {

			data := `
kind: Service
bar: foo
---
kind: Service
new: bar
---
kind: NOEXISTS
noway: togo
`

			tempdir, err := ioutil.TempDir("", "foo")
			Expect(err).ToNot(HaveOccurred())

			defer os.RemoveAll(tempdir)

			err = Generate(bytes.NewBufferString(data), tempdir)
			Expect(err).ToNot(HaveOccurred())

			serviceFile := filepath.Join(tempdir, "service.yaml")
			unknownFile := filepath.Join(tempdir, "other.yaml")
			Expect(serviceFile).To(BeARegularFile())
			Expect(unknownFile).To(BeARegularFile())

			content, err := ioutil.ReadFile(unknownFile)
			Expect(err).ToNot(HaveOccurred())

			serviceContent, err := ioutil.ReadFile(serviceFile)
			Expect(err).ToNot(HaveOccurred())

			Expect(string(content)).To(ContainSubstring("noway: togo"))
			Expect(string(serviceContent)).To(ContainSubstring("bar: foo"))

			Expect(string(serviceContent)).ToNot(ContainSubstring("noway: togo"))
			Expect(string(content)).ToNot(ContainSubstring("bar: foo"))

			Expect(string(serviceContent)).To(Equal(`kind: Service
bar: foo
---
kind: Service
new: bar
`), string(serviceContent))
		})

	})

	Context("Transforming", func() {
		It("replaces known helm templating values", func() {

			data := `kind: Service
metadata:
  name: some-metrics-service-name
  namespace: "foo"`

			tempdir, err := ioutil.TempDir("", "foo")
			Expect(err).ToNot(HaveOccurred())

			defer os.RemoveAll(tempdir)

			err = Generate(bytes.NewBufferString(data), tempdir, WithMutators(helm.Mutator))
			Expect(err).ToNot(HaveOccurred())

			serviceFile := filepath.Join(tempdir, "service.yaml")
			Expect(serviceFile).To(BeARegularFile())

			serviceContent, err := ioutil.ReadFile(serviceFile)
			Expect(err).ToNot(HaveOccurred())

			Expect(string(serviceContent)).To(ContainSubstring(`namespace: '{{.Release.Namespace}}'`), string(serviceContent))
			Expect(string(serviceContent)).To(ContainSubstring(`name: '{{ include "helm-chart.fullname" . }}-metrics-service'`), string(serviceContent))
		})
	})
})
