Errors
======

[![GoDoc](http://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/aodin/errors)

A robust errors type for Go.

### Quickstart

```go
package main

import (
    "encoding/json"
    "encoding/xml"
    "fmt"
    "log"
    "net/http"

    "github.com/aodin/errors"
)

func main() {
    err := errors.New()
    err.Code = http.StatusNotFound
    err.AddMeta("Not Found")
    err.SetField("ID", "Missing ID")

    // String
    fmt.Printf("%s\n", err.Error())
    // 404: Not Found; Missing ID (ID)

    // JSON
    jsonBytes, jsonErr := json.Marshal(err)
    if jsonErr != nil {
        log.Fatal(jsonErr)
    }
    fmt.Printf("%s\n", jsonBytes)
    // {"code":404,"meta":["Not Found"],"fields":{"ID":"Missing ID"}}

    // XML
    xmlBytes, xmlErr := xml.Marshal(err)
    if xmlErr != nil {
        log.Fatal(xmlErr)
    }
    fmt.Printf("%s\n", xmlBytes)
    // <Error><Code>404</Code><Metas><Meta>Not Found</Meta></Metas><Fields><ID>Missing ID</ID></Fields></Error>
}
```

Happy hacking!

aodin, 2016
