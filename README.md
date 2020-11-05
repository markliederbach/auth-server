# Auth Server <!-- omit in toc -->
[![Docker Image Version (tag latest semver)](https://img.shields.io/docker/v/markliederbach/auth-server/latest?label=docker%20image)](https://hub.docker.com/r/markliederbach/auth-server)

A home-grown JWT token provider built on [Gin](https://github.com/gin-gonic/gin). This project is highly experimental, and is not currently intended for anything more than a PoC and learning exercise.

- [Getting Started](#getting-started)
  - [Environment Variables](#environment-variables)
    - [`ACCESS_TOKEN_SECRET`/`REFRESH_TOKEN_SECRET` (required)](#access_token_secretrefresh_token_secret-required)
    - [`LOG_LEVEL` (optional)](#log_level-optional)
    - [`ACCESS_TOKEN_EXPIRE` (optional)](#access_token_expire-optional)
    - [`REFRESH_TOKEN_EXPIRE` (optional)](#refresh_token_expire-optional)
    - [`ISSUER` (optional)](#issuer-optional)
  - [Docker Container](#docker-container)
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

# Getting Started

## Environment Variables
There are several required and optional variables that can be passed to a running container to configure it.

### `ACCESS_TOKEN_SECRET`/`REFRESH_TOKEN_SECRET` (required)
These are secrets used to sign/verify all access/refresh tokens. Both are required.

If you need to generate a new one, here's a quick command:
```bash
openssl rand -hex 64
```

Be sure to make a unique one for each secret.

### `LOG_LEVEL` (optional)
Defaults to `INFO`. Options include `TRACE`, `DEBUG`, `INFO`, `WARN`, and `FATAL`.

### `ACCESS_TOKEN_EXPIRE` (optional)
Sets how long an individual access token should be valid. For available duration formats, please see [here](https://golang.org/pkg/time/#ParseDuration). Defaults to `15s` (which is admittedly very short).

### `REFRESH_TOKEN_EXPIRE` (optional)
Sets how long a refresh token should be valid. For available duration formats, please see [here](https://golang.org/pkg/time/#ParseDuration). Defaults to `1m` (which is admittedly very short).

### `ISSUER` (optional)
The label given to all tokens for the `iss` field. Defaults to `markliederbach/auth-service`.

## Docker Container
The recommended way to run this server is via Docker.

Once you have these environment variables set up (for example, in a `.env` file somewhere), you can run the following command to serve the application.

```bash
docker run --env-file .env -it -p 8080:8080/tcp --rm markliederbach/auth-server:latest
```

Of course, you can replace the host port with whatever you like with `-p 8080:<other port>/tcp`.

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