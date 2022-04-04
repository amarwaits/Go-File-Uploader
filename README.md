# Go File uploader

A simple HTTP server that serves a HTML form Upload image file to and store in server using Go.

## Prerequisites

* Valid Golang installation

## Getting started

Start the application.

```bash
export token=<your_preferred_token_here>    // export token=dh3bnj24v242nj
go run main.go
```

Hit `http://localhost:8080/upload` and you'll see a form to upload a file. After selecting a file and clicking upload, the file should be created in your local file-system.

## Implementation

* [Go File Uploader](https://github.com/amargc/brankas-go/blob/main/main.go)
