# File Server

File server provides a small server and a client to view and edit content on a file system.

Currently, only local file system is supported.

## How to Run

1. Docker Image

This application is containerized so a simple way is to build and use the docker images

You can build the docker image for server using:

> Optionally, you can set the `SERVER_IMAGE_NAME` or `IMAGE_REPO` to change the image.

```bash
$$ make build-server-image
```

In order to start the server, serving only the `test` directory you can use:

```bash
$$ make run-server-image
```

and this will listen on port 6000

Now you can run the client image viewing an arbitrary path using:

```bash
$$ make run-client-image QUERY=<path: for example: sub or sub/b.txt>
```

2. Binary

Make sure you have the golang installed.

### Build the server image using:

```bash
$$ make fs-server
```

Then you can run the server for an arbitrary path:

```bash
$$ ./bin/fs-server --root test
```

### Build the client image using:

```bash
$$ make fsc
```

Then you can run the client for an arbitrary path:

```bash
$$ ./bin/fsc --insecure <PATH>
```

You can optionally pass `--raw` to get the raw output.

## How to run test

You can simply run:

```
$$ make test
```

## Scripts for manual testing

Once you have the server running at port 6000, you can use the scripts in the `scripts` directory to
test out various different functions of the API, with a random file `c.txt`

## Repository Structure

There are two main packages, `filesystem` and the `fshttp`.

Since there are many static file hosting solutions such as nginx, the utility of an application like this would be in
the flexibility it can offer. Thus, an abstraction has been created to separate HTTP handlers from the file serving
logic.

`filesystem` package contains interfaces for any service that can provide file serving, this could be in-memory file
hosting, that can be used for testing or small short-lived application or it could be local file manager. In future,
we could provide S3 and Google cloud providers in order to support AWS and Google solutions using the very same HTTP
interface.

`fshttp` is the HTTP interface to the filesystem, it takes a filesystem.Editor and provides an HTTP interface to it.

Note that `fshttp` can simply be used in any other package, no dependency on other server applications has been added
so that anyone with any web framework could utilize this file serving utility.

### gRPC

One could simply provide gRPC capabilities to this project by adding yet another interface to it.
The adventage of such service would be that we could, then, using a grpc-web-proxy provide an HTTP
backend as well. And the nicest part is that we could generate swagger documentation.

I originally planned to use the gRPC but then time didn't allow for this, so this remains as future development.

## Helm Chart

This project does have a helm chart that can be used to deploy this application in a Kubernetes cluster.

Checkout `values.yaml` to see what values you can use at the moment. A more improved version of this helm chart
would allow for more customization in the service and deployment security settings like setting the user space.
