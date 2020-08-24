# go-git-providers

## Deprecated

> :warning: The successor to this library is https://github.com/fluxcd/go-git-providers

`go-git-providers` is a library that contains Go clients for different Git providers like Github.

At the moment only Github operations regarding deploy keys are supported.

**Warning**: The APIs contained in this repo are very unstable as it is very
early to define a good interface that will work well for most providers.

## Integration tests

To run the integration tests set the `GITHUB_TOKEN` environment variable, pick a Github repo to
run the tests on and run

```
make TEST_REPO=git@github.com:<USER>/go-git-providers.git integration-test
```
