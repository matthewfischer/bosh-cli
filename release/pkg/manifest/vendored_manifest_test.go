package manifest_test

import (
	"errors"

	fakesys "github.com/cloudfoundry/bosh-utils/system/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/cloudfoundry/bosh-cli/release/pkg/manifest"
)

var _ = Describe("NewVendoredManifestFromPath", func() {
	var (
		fs *fakesys.FakeFileSystem
	)

	BeforeEach(func() {
		fs = fakesys.NewFakeFileSystem()
	})

	It("parses pkg manifest successfully", func() {
		contents := `---
name: name
fingerprint: fp
dependencies:
- pkg1
- pkg2
`

		fs.WriteFileString("/path", contents)

		manifest, err := NewVendoredManifestFromPath("/path", fs)
		Expect(err).ToNot(HaveOccurred())
		Expect(manifest).To(Equal(VendoredManifest{
			Name:         "name",
			Fingerprint:  "fp",
			Dependencies: []string{"pkg1", "pkg2"},
		}))
	})

	It("returns error if manifest is not valid yaml", func() {
		fs.WriteFileString("/path", "-")

		_, err := NewVendoredManifestFromPath("/path", fs)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("line 1"))
	})

	It("returns error if manifest cannot be read", func() {
		fs.WriteFileString("/path", "-")
		fs.ReadFileError = errors.New("fake-err")

		_, err := NewVendoredManifestFromPath("/path", fs)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("fake-err"))
	})
})

var _ = Describe("VendoredManifest", func() {
	Describe("AsBytes", func() {
		It("returns serializes manifest", func() {
			bytes, err := VendoredManifest{
				Name: "name", Fingerprint: "fp", Dependencies: []string{"pkg1", "pkg2"}}.AsBytes()
			Expect(err).ToNot(HaveOccurred())
			Expect(string(bytes)).To(Equal(`name: name
fingerprint: fp
dependencies:
- pkg1
- pkg2
`))
		})

		It("returns error if name is empty", func() {
			_, err := VendoredManifest{}.AsBytes()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("Expected non-empty package name"))
		})

		It("returns error if fingerprint is empty", func() {
			_, err := VendoredManifest{Name: "name"}.AsBytes()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("Expected non-empty package fingerprint"))
		})
	})
})
