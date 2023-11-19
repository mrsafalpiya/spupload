# spupload

A simple file upload/host service with built in image optimization.

## Setup

### 1. Environment variables

- Copy `.env.example` to `.env`
- Edit the `.env` file
    - Replace `UPLOADS_DIR` (or `HOST_UPLOADS_DIR` if you are going to use docker) with the path you want the uploaded paths to be stored at.
    - (If you are using in production in an actual VPS) Replace `HOST` with the URL of the website.

### 2. Building binaries

#### i. Using docker (recommended)

```sh
$ docker-compose up
```

#### ii. Build it yourself

Make sure following dependencies are fulfilled:

- A proper Golang configuration
- `libwebp`

```sh
$ go mod download
$ go build -o ./out/spupload .
$ ./out/spupload
```

## License

GPLv3. See [COPYING](COPYING).
