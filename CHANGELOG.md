# Change Log
All notable changes to this project will be documented in this file.
This project adheres to [Semantic Versioning](http://semver.org/).

## [Unreleased]

## [1.0.0] - 2016-11-03
### Added
- `embedded` containers driver using cAdvisor under the hood.
- versioned configuration files with `version` the now-required field.
- `logging.fields` configuration value support to add custom log fields.
- `inmemory` storage driver support.
- `etcd` storage driver support.
- `consul` storage driver support.
- `http` and `http.cors` configuration sections to control HTTP API.
- basic hook CRUD support through the HTTP API.
- `slack+json` and `json` hook body format.
- `equal`, `eq`, `not_equal`, `ne`, and `match` condition operators.
