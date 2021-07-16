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

## Usage

```console
Archive web pages to disk.  When <url1> is www.google.com,
this tool will download the page and save it as www.google.com.html.

Usage:
  webdl [options] <url1> <url2> ...

Options:
  -debug
        show debug log
  -metadata
        show metadata
```

### Running `webdl` on console

Executing the following

```
$ ./webdl https://www.google.com https://twitter.com https://github.com
```

downloads the contents of each web pages into separate HTML files.

```
$ ls -l
github.com.html
twitter.com.html
www.google.com.html
```

### Using Docker

Using the Docker image created above,

```console
$ docker run --rm -v $(pwd):/tmp shihanng/webdl:0.1 https://www.google.com https://twitter.com https://github.com
```

## Unit-tests

Unit-tests is still a work in progress. Use `go test` to run tests.

```console
$ go test ./...
```
