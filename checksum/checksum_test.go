package checksum_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"fmt"

	. "github.com/cloudfoundry/bosh-utils/checksum"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
	fakesys "github.com/cloudfoundry/bosh-utils/system/fakes"
)

var _ = Describe("checksum", func() {
	Describe("#Verify", func() {
		It("verifies the algo and sum are matching", func() {
			expectedChecksum := NewChecksum("sha1", "07e1306432667f916639d47481edc4f2ca456454")
			actualChecksum := NewChecksum("sha1", "07e1306432667f916639d47481edc4f2ca456454")

			err := expectedChecksum.Verify(actualChecksum)
			Expect(err).ToNot(HaveOccurred())
		})

		Context("mismatching algorithm, matching checksum", func() {
			It("errors", func() {
				expectedChecksum := NewChecksum("sha1", "07e1306432667f916639d47481edc4f2ca456454")
				actualChecksum := NewChecksum("sha256", "07e1306432667f916639d47481edc4f2ca456454")

				err := expectedChecksum.Verify(actualChecksum)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal(`Expected sha1 algorithm but received sha256`))
			})
		})

		Context("matching algorithm, mismatching checksum", func() {
			It("errors", func() {
				expectedChecksum := NewChecksum("sha1", "07e1306432667f916639d47481edc4f2ca456454")
				actualChecksum := NewChecksum("sha1", "b1e66f505465c28d705cf587b041a6506cfe749f")

				err := expectedChecksum.Verify(actualChecksum)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal(`Expected sha1 checksum "07e1306432667f916639d47481edc4f2ca456454" but received "b1e66f505465c28d705cf587b041a6506cfe749f"`))
			})
		})
	})

	Describe("ChecksumFactory", func() {
		var (
			factory ChecksumFactory
			fs      boshsys.FileSystem
		)

		BeforeEach(func() {
			fs = fakesys.NewFakeFileSystem()
			factory = NewChecksumFactory(fs)
		})

		Describe("CreateFromFile", func() {
			const (
				filePath     = "/file.txt"
				fileContents = "something different"
			)

			BeforeEach(func() {
				fs.WriteFileString(filePath, fileContents)
			})

			Context("sha1", func() {
				It("opens a file and returns a checksum", func() {
					expectedChecksum, err := factory.CreateFromFile(filePath, "sha1")
					Expect(err).ToNot(HaveOccurred())
					Expect(expectedChecksum.Checksum()).To(Equal("da7102c07515effc353226eac2be923c916c5c94"))
				})
			})

			Context("sha256", func() {
				It("opens a file and returns a checksum", func() {
					expectedChecksum, err := factory.CreateFromFile(filePath, "sha256")
					Expect(err).ToNot(HaveOccurred())
					Expect(expectedChecksum.Checksum()).To(Equal("73af606b33433fa3a699134b39d5f6bce1ab4a6d9ca3263d3300f31fc5776b12"))
				})
			})

			Context("sha512", func() {
				It("opens a file and returns a checksum", func() {
					expectedChecksum, err := factory.CreateFromFile(filePath, "sha512")
					Expect(err).ToNot(HaveOccurred())
					Expect(expectedChecksum.Checksum()).To(Equal("25b38e5cf4069979d4de934ed6cde40eceec1f7100fc2a5fc38d3569456ab2b7e191bbf5a78b533df94a77fcd48b8cb025a4b5db20720d1ac36ecd9af0c8989a"))
				})
			})
		})
	})

	Describe("ParseString", func() {
		Describe("sha1", func() {
			It("creates a Checksum", func() {
				checksum, err := ParseString("sha1:07e1306432667f916639d47481edc4f2ca456454")
				Expect(err).ToNot(HaveOccurred())
				Expect(checksum.Algorithm()).To(Equal("sha1"))
				Expect(checksum.Checksum()).To(Equal("07e1306432667f916639d47481edc4f2ca456454"))
				Expect(checksum.String()).To(Equal("sha1:07e1306432667f916639d47481edc4f2ca456454"))
			})
		})

		Describe("sha256", func() {
			It("creates a Checksum", func() {
				checksum, err := ParseString("sha256:b1e66f505465c28d705cf587b041a6506cfe749f7aa4159d8a3f45cc53f1fb23")
				Expect(err).ToNot(HaveOccurred())
				Expect(checksum.Algorithm()).To(Equal("sha256"))
				Expect(checksum.Checksum()).To(Equal("b1e66f505465c28d705cf587b041a6506cfe749f7aa4159d8a3f45cc53f1fb23"))
				Expect(checksum.String()).To(Equal("sha256:b1e66f505465c28d705cf587b041a6506cfe749f7aa4159d8a3f45cc53f1fb23"))
			})
		})

		Describe("sha512", func() {
			It("creates a Checksum", func() {
				checksum, err := ParseString("sha512:6f06a0c6c3827d827145b077cd8c8b7a15c75eb2bed809569296e6502ef0872c8e7ef91307a6994fcd2be235d3c41e09bfe1b6023df45697d88111df4349d64a")
				Expect(err).ToNot(HaveOccurred())
				Expect(checksum.Algorithm()).To(Equal("sha512"))
				Expect(checksum.Checksum()).To(Equal("6f06a0c6c3827d827145b077cd8c8b7a15c75eb2bed809569296e6502ef0872c8e7ef91307a6994fcd2be235d3c41e09bfe1b6023df45697d88111df4349d64a"))
				Expect(checksum.String()).To(Equal("sha512:6f06a0c6c3827d827145b077cd8c8b7a15c75eb2bed809569296e6502ef0872c8e7ef91307a6994fcd2be235d3c41e09bfe1b6023df45697d88111df4349d64a"))
			})
		})

		Describe("default", func() {
			It("creates a sha1 Checksum", func() {
				checksum, err := ParseString("07e1306432667f916639d47481edc4f2ca456454")
				Expect(err).ToNot(HaveOccurred())
				Expect(checksum.Algorithm()).To(Equal("sha1"))
				Expect(checksum.Checksum()).To(Equal("07e1306432667f916639d47481edc4f2ca456454"))
				Expect(checksum.String()).To(Equal("sha1:07e1306432667f916639d47481edc4f2ca456454"))
			})
		})

		Describe("unrecognized", func() {
			It("errors", func() {
				_, err := ParseString("unrecognized:something")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("Unrecognized checksum algorithm: unrecognized"))
			})
		})
	})

	Describe("CreateHashFromAlgorithm", func() {
		data := []byte("the checksum of c1oudc0w is deterministic")

		Describe("sha1", func() {
			It("hashes", func() {
				hash, err := CreateHashFromAlgorithm("sha1")
				Expect(err).ToNot(HaveOccurred())

				hash.Write(data)
				Expect(fmt.Sprintf("%x", hash.Sum(nil))).To(Equal("07e1306432667f916639d47481edc4f2ca456454"))
			})
		})

		Describe("sha256", func() {
			It("hashes", func() {
				hash, err := CreateHashFromAlgorithm("sha256")
				Expect(err).ToNot(HaveOccurred())

				hash.Write(data)
				Expect(fmt.Sprintf("%x", hash.Sum(nil))).To(Equal("b1e66f505465c28d705cf587b041a6506cfe749f7aa4159d8a3f45cc53f1fb23"))
			})
		})

		Describe("sha512", func() {
			It("hashes", func() {
				hash, err := CreateHashFromAlgorithm("sha512")
				Expect(err).ToNot(HaveOccurred())

				hash.Write(data)
				Expect(fmt.Sprintf("%x", hash.Sum(nil))).To(Equal("6f06a0c6c3827d827145b077cd8c8b7a15c75eb2bed809569296e6502ef0872c8e7ef91307a6994fcd2be235d3c41e09bfe1b6023df45697d88111df4349d64a"))
			})
		})

		Describe("unrecognized", func() {
			It("errors", func() {
				_, err := CreateHashFromAlgorithm("unrecognized")
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
