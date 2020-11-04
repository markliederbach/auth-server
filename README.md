# Auth Server <!-- omit in toc -->

A home-grown JWT token provider built on [Gin](https://github.com/gin-gonic/gin). This project is highly experimental, and is not currently intended for anything more than a PoC and learning exercise.

- [Development](#development)
  - [Prerequisites](#prerequisites)
  - [Formatting](#formatting)
  - [Testing](#testing)
    - [Linting](#linting)
    - [Unit tests](#unit-tests)
  - [Building](#building)
    - [Local](#local)
    - [Remote](#remote)
  - [Release](#release)

# Development
## Prerequisites
The following must be installed for all other setup to work more easily.
- [`task`](https://taskfile.dev/#/installation)

After clonging the repo, simply run `task deps`. This will install [`pre-commit`](https://pre-commit.com/), [`goreleaser`](https://goreleaser.com/intro/), and other development dependencies. Additionally, it will bootstrap the needed commit hooks.

## Formatting
At any point, you can run the following to format your Go code, so the commit hooks and CI will pass.
```bash
task fmt
```

## Testing
### Linting
```bash
task lint
```
### Unit tests
There are currently no tests, as this is still in a PoC phase. However, the command to run tests exists as:
```
task test
```

## Building
### Local
To build artifacts for supported platforms, you can run the following:
```bash
task build
```
This will create the executables under `dist/`.

### Remote
The repo is also configured to build on push to the default branch, and upload artifacts for download.

## Release
This project has an automated release process through Github Actions with the help of [GoReleaser](https://goreleaser.com/intro/). When you are ready to create a new tagged release, do the following.

```bash
# Update `v0.1.0` to whatever the new version will be
git tag -a v0.1.0 -m "My new version" && git push --tags
```

This will kickstart the Github Action for releases. This workflow builds the binary for various
architectures/platforms, zips the archives, and uploads to a new draft release based on the tag.

Once you are satisfied with the release, and have added additional context to its description, you can publish the release.