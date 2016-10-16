# cSense

[![Unlicense](https://img.shields.io/badge/license-Unlicense-blue.svg?style=flat)](License) [![Report Card](https://goreportcard.com/badge/github.com/danielkrainas/csense)](goreportcard)

![cSense logo](https://github.com/danielkrainas/csense/blob/master/docs/logo/csense-logo.png)

cSense (Container Sense) allows you to create web hooks for when a container is created or exits. A web hook can have selectors for things like image name and metadata labels.

## Installation

> $ go get github.com/danielkrainas/csense

## Usage

> $ csense agent <config_path>

**Example** - with development config:

> $ csense agent ./config.dev.yml

## Configuration

A configuration file is *required* for cSense but environment variables can be used to override configuration. A configuration file can be specified as a parameter or with the `CSENSE_CONFIG_PATH` environment variable. 

All configuration environment variables are prefixed by `CSENSE_` and the paths are separated by an underscore(`_`). Some examples:

- `CSENSE_LOGGING_LEVEL=warn`
- `CSENSE_HTTP_ADDR=localhost:2345`
- `CSENSE_STORAGE_MOCK=true`
- `CSENSE_STORAGE_CONSUL_PARAM1=val`

A development configuration file is included: `/config.dev.yml` and a `/config.local.yml` has already been added to gitignore to be used for local testing or development.

```yaml
# configuration schema version number, only `0.1`
version: 0.1

# log stuff
logging:
  # minimum event level to log: `error`, `warn`, `info`, or `debug`
  level: 'debug'
  # log output format: `text` or `json`
  formatter: 'text'
  # custom fields to be added and displayed in the log
  fields:
    customfield1: 'value'

# http server stuff
http:
  # host:port address for the server to listen on
  addr: ':9181'
  # http host
  host: 'localhost'

  # CORS stuff
  cors:
    # origins to allow
    origins: ['http://localhost:5555']
    # methods to allow
    methods: ['GET','POST','OPTIONS','DELETE','CONNECT']
    # headers to allow
    headers: ['*']

# storage driver and parameters
storage:
  consul:
    param1: 'val'

# the mock driver has no parameters so it can be declared as a string
storage: 'mock'
```

`storage` only allows specification of *one* driver per configuration. Any additional ones will cause a validation error when the application starts.

## Development

This is information related to developing cSense itself.

### Tools

These are tools used for development.

#### Vendoring

`gvt` is used to track and vendor dependencies, install with:

> $ go get github.com/FiloSottile/gvt

For details on usage, see [the project's github](https://github.com/FiloSottile/gvt).

### Building

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

## License

[Unlicense](http://unlicense.org/UNLICENSE). This is a Public Domain work. 

[![Public Domain](https://licensebuttons.net/p/mark/1.0/88x31.png)](http://questioncopyright.org/promise)

> ["Make art not law"](http://questioncopyright.org/make_art_not_law_interview) -Nina Paley