![GitHub repo file count](https://img.shields.io/github/directory-file-count/zagart47/imageserver)
![GitHub go.mod Go version (subdirectory of monorepo)](https://img.shields.io/github/go-mod/go-version/zagart47/imageserver?filename=go.mod)
![GitHub last commit](https://img.shields.io/github/last-commit/zagart47/imageserver)
# Image Server.

The service knows how to upload, count and download files on demand.
To use the service, you need to use the client, the link to it is presented below.
The service is written as part of a test assignment.
The service uses two ports. Port ```80``` is used to upload and download files. Port ```81``` is used to view the file list.

## Installation

Use the git cli to clone repository.

```bash
git clone https://github.com/zagart47/imageserver.git
```

## Usage
### Start server
```bash
cd imageserver
go run cmd/app/main.go
```

# Attention! You need use the client to work with this server!!!

## Client here >>>>>[github.com/zagart47/imageclient](https://github.com/zagart47/imageclient)<<<<<


