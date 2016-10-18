![cSense logo](https://github.com/danielkrainas/csense/blob/master/docs/logo/csense-logo.png)

# Build and Run Tests

## Building 

One caveat with building currently is that because of the cAdvisor dependency for the containers driver, `cgo` *cannot* be disabled; the build will fail. So no `CGO_ENABLED=0` builds.

#### Dev/Local build

Use `go` and build from the root of the project:

> $ go build

Please note the version number displayed will be the value of `main.DEFAULT_VERSION`

#### Local versioned build

Use `make` to create a versioned build:

> $ make compile

The default version is a semver-compatible string made up of the contents of the `/VERSION` file and the short form of the current git hash (e.g: `1.0.0-c63076f`). To override this default version, you may use the `BUILD_VERSION` environment variable to set it manually:

> $ BUILD_VERSION=7.7.7-lucky make compile

#### Dist build

This is primarily meant to be used when building the docker image. Distribution builds are versioned like the local versioned builds and the build specifically targets `linux`

> $ make dist

#### Docker Image

Building a Docker Image is a two-step process because of the CGO requirement and the desire to keep a small image size. First we build the distribution binary:

> $ make dist

And then we can make the image:

> $ make image

The default image repo used is that of the Makefile's `DOCKER_REPO` variable. The image tag is the `BUILD_VERSION` variable and can be overridden as noted in the *"Local versioned build"* section above.

### Testing

Use `make` to run tests:

> $ make test

You can also use `go test` directly for any package without additional bootstrapping:

> $ go test ./api/