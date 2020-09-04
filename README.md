# delrange

delrange is a static analysis tool which detects delete function is called with a value different from range key.

![test_and_lint](https://github.com/p1ass/delrange/workflows/test_and_lint/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/p1ass/delrange)](https://goreportcard.com/report/github.com/p1ass/delrange)


## Features

- Detect delete function is called with a value different from range key.



### Example

```go
package main

import "fmt"

func main() {
	m := map[int]int{1: 1, 2: 2}
	for key, value := range m {
		delete(m, 1) // want "function is called with a value different from range key"
		delete(m, key)
		fmt.Println(key, value)
	}
}

```

## Installation

### go get

```shell script
GO111MODULE=off go get github.com/p1ass/delrange/cmd/delrange
```

## Usage

```shell script
go vet -vettool=`which delrange` ./...
``` 