package parser_test

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/spectrocloud/kubesplit/internal/parser"
)

var _ = Describe("Parser", func() {
	Context("Splitting elements", func() {
		It("Splits them correctly by type", func() {

			data := `
kind: Service
bar: foo
---
kind: NOEXISTS
noway: togo
`

			tempdir, err := ioutil.TempDir("", "foo")
			Expect(err).ToNot(HaveOccurred())

			defer os.RemoveAll(tempdir)

			Generate(bytes.NewBufferString(data), tempdir)

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
		})

	})

})
