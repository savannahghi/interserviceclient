[![Build Status](https://travis-ci.com/savannahghi/interserviceclient.svg?branch=main)](https://travis-ci.com/savannahghi/interserviceclient)
[![Maintained](https://img.shields.io/badge/Maintained-Actively-informational.svg?style=for-the-badge)](https://shields.io/)
# Interservice Library
interserviceclient  play a major role in implementing Microservices Inter-service communication using either HTTP protocol or asynchronous message patterns.

### Installing it
interserviceclient is compatible with modern Go releases in module mode, with Go installed:

```bash
go get -u github.com/savannahghi/interserviceclient

```
will resolve and add the package to the current development module, along with its dependencies.

Alternatively the same can be achieved if you use import in a package:

```go
import "github.com/savannahghi/interserviceclient"

```
and run `go get` without parameters.

The package name is `interserviceclient`


### Developing

The default branch library is `main`

We try to follow semantic versioning ( <https://semver.org/> ). For that reason,
every major, minor and point release should be _tagged_.

```
git tag -m "v0.0.1" "v0.0.1"
git push --tags
```

Continuous integration tests *must* pass on Travis CI. Our coverage threshold
is 90% i.e you *must* keep coverage above 90%.


## Environment variables

In order to run tests, you need to have an `env.sh` file similar to this one:

```bash
# Application settings
export DEBUG=true
export IS_RUNNING_TESTS=true
export SENTRY_DSN=<a Sentry Data Source Name>

# Google Cloud credentials
export GOOGLE_APPLICATION_CREDENTIALS="<path to a service account JSON file"
export GOOGLE_CLOUD_PROJECT="Google Cloud project id"
export FIREBASE_WEB_API_KEY="<a web API key that corresponds to the project named above>"
export FIREBASE_DYNAMIC_LINKS_DOMAIN=https://bwlci.page.link
# Firestore documents root collection suffix
export ROOT_COLLECTION_SUFFIX="testing"

```

This file *must not* be committed to version control.

It is important to _export_ the environment variables. If they are not exported,
they will not be visible to child processes e.g `go test ./...`.

These environment variables should also be set up on Travis CI environment variable section.

## Contributing ##
I would like to cover the entire GitHub API and contributions are of course always welcome. The
calling pattern is pretty well established, so adding new methods is relatively
straightforward. See [`CONTRIBUTING.md`](CONTRIBUTING.md) for details.

## Versioning ##

In general, interserviceclient follows [semver](https://semver.org/) as closely as we
can for tagging releases of the package. For self-contained libraries, the
application of semantic versioning is relatively straightforward and generally
understood. We've adopted the following
versioning policy:

* We increment the **major version** with any incompatible change to
	non-preview functionality, including changes to the exported Go API surface
	or behavior of the API.
* We increment the **minor version** with any backwards-compatible changes to
	functionality, as well as any changes to preview functionality in the GitHub
	API. GitHub makes no guarantee about the stability of preview functionality,
	so neither do we consider it a stable part of the go-github API.
* We increment the **patch version** with any backwards-compatible bug fixes.

## License ##

This library is distributed under the MIT license found in the [LICENSE](./LICENSE)
file.