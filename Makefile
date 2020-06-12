.PHONY: test
integration-test: check-github-token
	go test ./pkg/integration... -test.v -ginkgo.v -repo $(TEST_REPO)


check-github-token:
ifndef GITHUB_TOKEN
	$(error GITHUB_TOKEN is undefined)
endif