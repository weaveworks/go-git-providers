package integration

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"flag"
	"fmt"
	"strings"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"golang.org/x/crypto/ssh"

	"github.com/weaveworks/go-git-providers/pkg/key"
	"github.com/weaveworks/go-git-providers/pkg/providers"
)

var testRepo string

func init() {
	// Call testing.Init() prior to tests.NewParams(), as otherwise -test.* will not be recognised. See also: https://golang.org/doc/go1.13#testing
	testing.Init()
	flag.StringVar(&testRepo, "repo", "", "The ssh url to the repo used for integration testing")
}

var _ = Describe("Github Provider", func() {
	BeforeSuite(func() {
		flag.Parse()
		if testRepo == "" {
			Fail("Missing --repo argument")
		}
	})

	It("should authorize and delete an SSH key", func() {
		provider, err := providers.GetProvider(testRepo)
		Expect(err).ToNot(HaveOccurred())

		testKey := createSSHKey("integ-test-1", false)

		// Authorize a deploy key
		err = provider.AuthorizeSSHKey(context.Background(), testKey)
		Expect(err).ToNot(HaveOccurred())

		// Check it was uploaded
		checkKeyUploaded(provider, testKey)

		// Delete the key
		err = provider.DeleteSSHKey(context.Background(), testKey.Title)
		Expect(err).ToNot(HaveOccurred())

		// Check it was deleted
		checkKeyDeleted(provider, testKey)
	})

	It("should authorize an SSH key as read write", func() {
		provider, err := providers.GetProvider(testRepo)
		Expect(err).ToNot(HaveOccurred())

		// Upload deploy key
		testKey := createSSHKey("integ-test-2", true)
		err = provider.AuthorizeSSHKey(context.Background(), testKey)
		Expect(err).ToNot(HaveOccurred())

		// Check key was uploaded
		checkKeyUploaded(provider, testKey)

		// Delete key
		err = provider.DeleteSSHKey(context.Background(), testKey.Title)
		Expect(err).ToNot(HaveOccurred())

		// Check key was deleted
		checkKeyDeleted(provider, testKey)
	})
})

func createSSHKey(title string, readOnly bool) key.SSHKey {
	publicKey, err := generatePublicKey()
	Expect(err).ToNot(HaveOccurred())

	return key.SSHKey{
		Title:    fmt.Sprintf("%s-%s", title, time.Now().Format("2006-01-02T15:04:05")),
		Key:      strings.TrimSpace(publicKey),
		ReadOnly: readOnly,
	}
}

func generatePublicKey() (string, error) {
	// Private Key generation
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return "", err
	}

	publicRsaKey, err := ssh.NewPublicKey(&privateKey.PublicKey)
	if err != nil {
		return "", err
	}

	pubKeyBytes := ssh.MarshalAuthorizedKey(publicRsaKey)

	return string(pubKeyBytes), nil
}

func findKey(keys []key.SSHKey, title string) *key.SSHKey {
	for _, k := range keys {
		if k.Title == title {
			return &k
		}
	}
	return nil
}

func checkKeyUploaded(provider providers.Provider, testKey key.SSHKey) {
	uploadedKeys, err := provider.ListKeys(context.Background())
	Expect(err).ToNot(HaveOccurred())

	uploadedKey := findKey(uploadedKeys, testKey.Title)
	Expect(uploadedKey).ToNot(BeNil())
	Expect(uploadedKey.Title).To(Equal(testKey.Title))
	Expect(uploadedKey.Key).To(Equal(testKey.Key))
	Expect(uploadedKey.ReadOnly).To(Equal(testKey.ReadOnly))
}

func checkKeyDeleted(provider providers.Provider, testKey key.SSHKey) {
	uploadedKeys, err := provider.ListKeys(context.Background())
	Expect(err).ToNot(HaveOccurred())

	uploadedKey := findKey(uploadedKeys, testKey.Title)
	Expect(uploadedKey).To(BeNil())
}