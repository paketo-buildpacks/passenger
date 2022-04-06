package passenger_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/paketo-buildpacks/passenger"
	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"
)

func testPassengerfileParser(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		path   string
		parser passenger.PassengerfileParser
	)

	it.Before(func() {
		file, err := os.CreateTemp("", "Passengerfile.json")
		Expect(err).NotTo(HaveOccurred())
		defer file.Close()

		path = file.Name()

		Expect(os.WriteFile(path, []byte(`{"port":4000}`), 0644)).To(Succeed())

		parser = passenger.NewPassengerfileParser()
	})

	it.After(func() {
		Expect(os.RemoveAll(path)).To(Succeed())
	})

	context("ParsePort", func() {
		it("returns the port from the Passengerfile", func() {
			passengerfile, err := parser.ParsePort(path, 1234)
			Expect(err).NotTo(HaveOccurred())

			Expect(passengerfile).To(Equal(4000))
		})

		context("when the Passengerfile does not exist", func() {
			it.Before(func() {
				Expect(os.Remove(path)).To(Succeed())
			})

			it("returns the provided default port", func() {
				port, err := parser.ParsePort(path, 1234)
				Expect(err).NotTo(HaveOccurred())

				Expect(port).To(Equal(1234))
			})
		})

		context("when the Passengerfile does not contain the port field", func() {
			it.Before(func() {
				Expect(os.WriteFile(path, []byte(`{"some-other-field":"some-other-value"}`), 0644)).To(Succeed())
			})

			it("returns the provided default port", func() {
				port, err := parser.ParsePort(path, 1234)
				Expect(err).NotTo(HaveOccurred())

				Expect(port).To(Equal(1234))
			})
		})

		context("failure cases", func() {
			context("when determining if the Passengerfile exists fails", func() {
				var (
					tempDir string
				)

				it.Before(func() {
					var err error
					tempDir, err = os.MkdirTemp("", "")
					Expect(err).NotTo(HaveOccurred())

					Expect(os.Chmod(tempDir, 0000)).To(Succeed())
				})

				it.After(func() {
					Expect(os.RemoveAll(tempDir)).To(Succeed())
				})

				it("returns an error", func() {
					_, err := parser.ParsePort(filepath.Join(tempDir, "some-file"), 1234)
					Expect(err).To(HaveOccurred())
					Expect(err).To(MatchError(ContainSubstring("failed to determine if Passengerfile exists:")))
					Expect(err).To(MatchError(ContainSubstring("some-file")))
				})
			})

			context("when the Passengerfile cannot be opened", func() {
				it.Before(func() {
					Expect(os.Chmod(path, 0000)).To(Succeed())
				})

				it("returns an error", func() {
					_, err := parser.ParsePort(path, 1234)
					Expect(err).To(HaveOccurred())
					Expect(err).To(MatchError(ContainSubstring("failed to read Passengerfile:")))
					Expect(err).To(MatchError(ContainSubstring("permission denied")))
				})
			})

			context("when the Passengerfile is malformed", func() {
				it.Before(func() {
					Expect(os.WriteFile(path, []byte(`{"port":"not-an-integer"}`), 0644)).To(Succeed())
				})

				it("returns an error", func() {
					_, err := parser.ParsePort(path, 1234)
					Expect(err).To(HaveOccurred())
					Expect(err).To(MatchError(ContainSubstring("failed to parse Passengerfile:")))
				})
			})
		})
	})
}
