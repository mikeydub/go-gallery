# go-gallery

## Pre

1. [Install Go](https://golang.org/doc/install)
2. [Install Docker](https://www.docker.com/products/docker-desktop)
3. [Install Docker-Compose](https://docs.docker.com/compose/install/)

## Clone and install deps

```bash
$ git clone git@github.com:gallery-so/go-gallery.git
$ cd go-gallery
$ go get -u -d ./...
```

## Setup (Mac)

```bash
$ go build -o ./bin/main ./cmd/server/main.go
```

This will generate a binary within `./bin/main`. To run the binary, simply:

```bash
$ ./bin/main
```

### Mongo and Redis

The app will connect to a local mongo instance and local redis instance by default. To spin it up, you can use the official docker containers:

```bash
$ docker-compose up -d
```

To remove running mongo and redis instance:

```bash
$ docker-compose down
```

### Healthcheck

Verify that the server is running by calling the `/v1/health` endpoint.

```bash
$ curl localhost:4000/v1/health
```

## Testing

Run all tests in current directory and all of its subdirectories:

```bash
$ go test ./...
```

Run all tests in subdirectory (e.g. /server):

```bash
$ go test ./server/...
```

Run a specific test by passing the name as an option:

```bash
go test -run {testName}
```

Add `-v` for detailed logs
