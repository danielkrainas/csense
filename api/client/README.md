# cSense API Client

Client library for the cSense API. 

Supported Endpoints:

- Hooks


## Installation

> $ go get github.com/danielkrainas/csense/api/client


## Usage

How to instantiate a new client:

```go
package main

import (
	"net/http"
	"github.com/danielkrainas/csense/api/client"
)

// http/https url of the csense service
const ENDPOINT = "http://localhost:9181"

func main() {
	// Create a new client
	c := client.New(ENDPOINT, http.DefaultClient)
}
```


## Example

A more detailed example can be found [here.](https://github.com/danielkrainas/csense/tree/master/api/client/example)

