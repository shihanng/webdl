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
Archive web pages to disk.  When <url1> is https://www.google.com,
this tool will download the page and save it as index.html
in www.google.com/ directory.

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
$ ./webdl https://developer.mozilla.org/en-US/ https://www.google.com
```

downloads the contents and assets of each web pages into separate HTML files and store them into separate directories based on their host name.

```console
├── developer.mozilla.org
│   ├── apple-touch-icon.0ea0fa02.png
│   ├── en-US_.html
│   ├── favicon-48x48.97046865.png
│   ├── manifest.56b1cedc.json
│   └── static
│       ├── css
│       │   └── main.20e8790b.chunk.css
│       ├── js
│       │   ├── 2.c6ebdd97.chunk.js
│       │   ├── ga.js
│       │   ├── main.f32a715f.chunk.js
│       │   └── runtime-main.6526a4ac.js
│       └── media
│           └── ZillaSlab-Bold.subset.0beac26b.woff2
└── www.google.com
    ├── images
    │   └── branding
    │       └── googlelogo
    │           └── 1x
    │               └── googlelogo_white_background_color_272x92dp.png
    └── index.html
```

Using a simple web server, we can view the downloaded HTML together with its assets. E.g.

```console
$ cd developer.mozilla.org
$ python -m http.server
```

and visit <http://localhost:8000/>.


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
