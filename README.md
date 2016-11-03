# cSense

[![License](https://img.shields.io/badge/license-Unlicense-blue.svg?style=flat)](UNLICENSE) [![Go Report Card](https://goreportcard.com/badge/github.com/danielkrainas/csense)](https://goreportcard.com/report/github.com/danielkrainas/csense) [![Docker Hub](https://img.shields.io/docker/pulls/dakr/csense.svg?style=flat)](https://hub.docker.com/r/dakr/csense/)

![cSense logo](https://github.com/danielkrainas/csense/blob/master/docs/logo/csense-logo.png)

cSense (Container Sense) allows you to subscribe to container events with web hooks. Hooks are registered with cSense via an HTTP API and may contain selectors, like image tag or container name, to limit the containers or events that the hook should be notified about.

## Installation

> $ go get github.com/danielkrainas/csense

## Usage

> $ csense [command] <config_path>

Most commands require a configuration path provided as an argument or in the `CSENSE_CONFIG_PATH` environment variable. 

### Agent mode

This is the primary mode for cSense. It hosts the HTTP API server and handles monitoring and notifying hooks of container events.

> $ csense agent <config_path>

**Example** - with the default config:

> $ csense agent ./config.default.yml

## Configuration

A configuration file is *required* for cSense but environment variables can be used to override configuration. A configuration file can be specified as a parameter or with the `CSENSE_CONFIG_PATH` environment variable. 

All configuration environment variables are prefixed by `CSENSE_` and the paths are separated by an underscore(`_`). Some examples:

- `CSENSE_LOGGING_LEVEL=warn`
- `CSENSE_HTTP_ADDR=localhost:2345`
- `CSENSE_STORAGE_INMEMORY=true`
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

# the in-memory driver has no parameters so it can be declared as a string
storage: 'inmemory'
```

`storage` only allows specification of *one* driver per configuration. Any additional ones will cause a validation error when the application starts.

## Bugs and Feedback

If you see a bug or have a suggestion, feel free to open an issue [here](https://github.com/danielkrainas/csense/issues).

## Contributions

PR's welcome! There are no strict style guidelines, just follow best practices and try to keep with the general look & feel of the code present. All submissions should atleast be `go fmt -s` and have a test to verify *(if applicable)*.

For details on how to extend and develop cSense, see the [dev documentation](docs/development/).

## License

[Unlicense](http://unlicense.org/UNLICENSE). This is a Public Domain work. 

[![Public Domain](https://licensebuttons.net/p/mark/1.0/88x31.png)](http://questioncopyright.org/promise)

> ["Make art not law"](http://questioncopyright.org/make_art_not_law_interview) -Nina Paley