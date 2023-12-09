FROM golang:alpine

RUN apk add --no-cache build-base libwebp-dev

# Set the Current Working Directory inside the container
WORKDIR /app/spupload

# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

# Build the Go app
RUN go build -o ./out/spupload .


# This container exposes port 8080 to the outside world
EXPOSE 5433

# Run the binary program produced by `go install`
CMD ["./out/spupload"]
