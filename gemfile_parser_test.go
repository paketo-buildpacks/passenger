package passenger_test

import (
	"os"
	"testing"

	"github.com/paketo-buildpacks/passenger"
	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"
)

func testGemfileParser(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		path   string
		parser passenger.GemfileParser
	)

	it.Before(func() {
		file, err := os.CreateTemp("", "Gemfile")
		Expect(err).NotTo(HaveOccurred())
		defer file.Close()

		path = file.Name()

		parser = passenger.NewGemfileParser()
	})

	it.After(func() {
		Expect(os.RemoveAll(path)).To(Succeed())
	})

	context("Parse", func() {
		context("when using passenger", func() {
			it("parses correctly", func() {
				const GEMFILE_CONTENTS = `source 'https://rubygems.org'
gem 'passenger'`

				Expect(os.WriteFile(path, []byte(GEMFILE_CONTENTS), 0644)).To(Succeed())

				hasPassenger, err := parser.Parse(path)
				Expect(err).NotTo(HaveOccurred())
				Expect(hasPassenger).To(Equal(true))
			})
		})

		context("when not using passenger", func() {
			it("parses correctly", func() {
				const GEMFILE_CONTENTS = `source 'https://rubygems.org'`

				Expect(os.WriteFile(path, []byte(GEMFILE_CONTENTS), 0644)).To(Succeed())

				hasPassenger, err := parser.Parse(path)
				Expect(err).NotTo(HaveOccurred())
				Expect(hasPassenger).To(Equal(false))
			})
		})

		context("when the Gemfile file does not exist", func() {
			it.Before(func() {
				Expect(os.Remove(path)).To(Succeed())
			})

			it("returns all false", func() {
				hasPassenger, err := parser.Parse(path)
				Expect(err).NotTo(HaveOccurred())
				Expect(hasPassenger).To(Equal(false))
			})
		})

		context("failure cases", func() {
			context("when the Gemfile cannot be opened", func() {
				it.Before(func() {
					Expect(os.Chmod(path, 0000)).To(Succeed())
				})

				it("returns an error", func() {
					_, err := parser.Parse(path)
					Expect(err).To(HaveOccurred())
					Expect(err).To(MatchError(ContainSubstring("failed to parse Gemfile:")))
					Expect(err).To(MatchError(ContainSubstring("permission denied")))
				})
			})
		})
	})
}
