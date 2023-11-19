# spupload

A simple file upload/host service with built in image optimization.

## Features

- Very low system requirements.
- Automatic image optimization to webp.

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

## Usage

Considering the server is hosted at `http://localhost:5433`, following endpoints are available:

### 1. Upload

```
POST http://localhost:5433/<folder_name>/<sub_folder_name>/...

Form Data:
file = File to upload (Required)
filename = Name of the file to upload as (Optional)
replace = "true" or "false" (Optional) (By default if we are uploading a file named `foo` and we already have a file named `foo`, the new file will be named `foo-1`. Setting this true will discard this feature and the file will be replaced)
disable-file-optimization = "true" or "false" (Optional) (Disable any file optimization, example: Image conversion to webp)
```

### 2. Receive uploaded file

```
GET http://localhost:5433/<folder_name>/<sub_folder_name>/.../<file_name>

Query params:
view = "detail" (Optional) (Get details about the file: Filename, Size, Modification Time, Filetype)
```

## License

GPLv3. See [COPYING](COPYING).
