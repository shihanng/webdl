# webdl

Simple CLI tool to archive web pages using [Colly](https://github.com/gocolly/colly).

## Build

Choose either "from source" or "using Docker."

### From source

[Go](https://golang.org/) is required for this step.
The following will give an executable called `webdl`.

```console
$ go build -o webdl
```

### Using Docker

We also provide Dockerfile to build Docker image of the CLI tool. [Docker](https://www.docker.com/) is required for this.
The following builds a Docker image `shihanng/webdl:0.1`.

```console
$ docker build -t shihanng/webdl:0.1 .
```

## Unit-tests

Unit-tests is still a work in progress. Use `go test` to run tests.

```console
$ go test ./...
```
