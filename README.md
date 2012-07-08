# Aspell library bindings for Go

GNU Aspell is a spell checking tool written in C/C++. This package provides simplified Aspell bindings for Go.
It uses UTF-8 by default and encapsulates some Aspell internals.

## Getting started

First make sure aspell library and headers are installed on your system. On Debian/Ubuntu you could install it this way:

```
sudo apt-get install aspell libaspell-dev
```

It you need some more dictionaries you can install them like this:

```
sudo apt-get install aspell-ua aspell-se
```

Then you can install the package using the Go tool:

```
go get github.com/trustmaster/go-aspell
```

## Usage

Here is a simple spell checker program using the aspell package:

```go
package main

import (
	"github.com/trustmaster/go-aspell"
	"fmt"
	"os"
	"strings"
)

func main() {
	// Get a word from cmd line arguments
	if len(os.Args) != 2 {
		fmt.Print("Usage: aspell_example word\n")
		return
	}
	word := os.Args[1]

	// Initialize the speller
	speller, err := aspell.NewSpeller(map[string]string{
		"lang": "en_US",
	})
	if err != nil {
		fmt.Errorf("Error: %s", err.Error())
		return
	}
	defer speller.Delete()

	// Check and suggest
	if speller.Check(word) {
		fmt.Print("OK\n")
	} else {
		fmt.Printf("Incorrect word, suggestions: %s\n", strings.Join(speller.Suggest(word), ", "))
	}
}
```

For more information see [aspell_test.go](https://github.com/trustmaster/go-aspell/blob/master/aspell_test.go) file and use the godoc tool:

```
godoc github.com/trustmaster/go-aspell
```
