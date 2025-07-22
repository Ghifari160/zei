<!-- markdownlint-configure-file {
    "hard_tab": {
        "ignore_code_languages": [ "go","golang" ]
    }
} -->
# Zei

Go HTTP Client library built atop the standard library's [net/http.Client].

## Install

``` shell
go get -u github.com/ghifari160/zei
```

## Usage

[Documentation]

### Quick Start

``` go
package main

import "github.com/ghifari160/zei"

func main() {
	client := zei.New(&zei.Config{})
	resp, err := client.Get("http://example.com/")
	if err != nil {
		// Handle error
	}
	defer resp.Body.Close()
	// ...
}
```

### Set user agent

``` go
client := zei.New(&zei.Config{UserAgent: "App/1.0"})
resp, err := client.Get("http://example.com/")
```

### Set timeout

``` go
client := zei.New(&zei.Config{Timeout: 1 * time.Minute})
resp, err := client.Get("http://example.com/")
```

### Set authorization

``` go
config := zei.Config{}
// Basic Authentication
config.SetBasicAuth("username", "password")
// Bearer Authentication
config.SetBearerAuth("token_value")

client := zei.New(&config)
resp, err := client.Get("http://example.com/")
```

[Documentation]: https://pkg.go.dev/github.com/ghifari160/zei
[net/http.Client]: https://pkg.go.dev/net/http#Client
