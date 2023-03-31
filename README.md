# gobump CLI tool

[![CI](https://github.com/danilvpetrov/gobump/actions/workflows/ci.yml/badge.svg)](https://github.com/danilvpetrov/gobump/actions/workflows/ci.yml)

## About

This tool allows managing the major version in the Go module paths. The module
path can be the path of the module itself or one of the module's dependencies

The following example upgrades the module path to `v2`. The command requires
running within a valid Go module directory.

```sh
gobump github.com/exampleorg/examplerepo/v2
```al

The following example downgrades the module path to `v0/v1`. The command
requires running within a valid Go module directory.

```sh
gobump github.com/exampleorg/examplerepo
```

## Installation

To install into `GOBIN` folder run the following command:

```sh
go install github.com/danilvpetrov/gobump/cmd/gobump@latest
```
