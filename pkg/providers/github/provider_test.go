package github

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Github Provider", func() {
	It("should parse simple urls correctly", func() {
		testURL := "git@github.com:weaveworks/go-git-provider"
		name, err := repoName(testURL)

		Expect(err).ToNot(HaveOccurred())
		Expect(name).To(Equal("go-git-provider"))

		owner, err := repoOwner(testURL)
		Expect(owner).To(Equal("weaveworks"))
	})

	It("should parse correctly when it ends in .git", func() {
		testURL := "git@github.com:weaveworks/go-git-provider.git"
		name, err := repoName(testURL)

		Expect(err).ToNot(HaveOccurred())
		Expect(name).To(Equal("go-git-provider"))

		owner, err := repoOwner(testURL)
		Expect(owner).To(Equal("weaveworks"))
	})

	It("should parse correctly when it starts with ssh://", func() {
		testURL := "ssh://git@github.com:weaveworks/go-git-provider"
		name, err := repoName(testURL)

		Expect(err).ToNot(HaveOccurred())
		Expect(name).To(Equal("go-git-provider"))

		owner, err := repoOwner(testURL)
		Expect(owner).To(Equal("weaveworks"))
	})

	It("should parse correctly when it starts with ssh:// and ends with .git", func() {
		testURL := "ssh://git@github.com:weaveworks/go-git-provider.git"
		name, err := repoName(testURL)

		Expect(err).ToNot(HaveOccurred())
		Expect(name).To(Equal("go-git-provider"))

		owner, err := repoOwner(testURL)
		Expect(owner).To(Equal("weaveworks"))
	})
})