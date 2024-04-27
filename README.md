
## g4fi

**g4f-interactive** - A simple command line client for [g4f](https://github.com/xtekky/gpt4free)

## Installation

```
git clone https://github.com/sheepla/g4fi.git
go install
```

## Usage

Using this tool requires the g4f server to be running and the ability to communicate with the server on the corresponding port number. If you are using Docker, you can start the container with the following:

```sh
docker run -p 8080:8080 -p 1337:1337 -p 7900:7900 --shm-size="2g"  hlohaus789/g4f:latest
```

Command line usage is below:

```
Usage: g4fi [--timeout TIMEOUT] [--server SERVER]

Options:
  --timeout TIMEOUT, -t TIMEOUT
                         Timeout seconds [default: 30]
  --server SERVER, -s SERVER
                         hostname and port of g4f API instance [default: localhost:8080]
  --help, -h             display this help and exit
```

When you launch the command, it initiates an interactive session. 
This session allows to input prompt messages and receive corresponding results. 
You can utilize line editing with key bindings similar to GNU ReadLine, to efficiently input your prompt messages.. 
To exit the interactive session, simply type `Ctrl-D`.

```
[you@your-computer]$ g4fi

> How to use net/http in Go?

The `net/http` package in Go is used for HTTP client and server implementations. Here's a basic example of how to use `net/http` to create a simple HTTP server and client.

**HTTP Server:**
```go
package main

import (
    "fmt"
    "net/http"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w

>
```

