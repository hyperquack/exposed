![Banner](assets/banner.png)

# Exposed ("Client")

Exposed is a Continuous Threat Exposure system that alerts when your attack surface changes.

## Usage

> A valid keychain file is required to get started

Simple to install:

```go
go install github.com/hyperquack/exposed@latest
```

Simple to use:

```zsh
exposed targets                         # list your target domains
exposed start example.com               # start monitoring a new target
exposed stop example.com                # stop monitoring a target
exposed notify https://hooks.slack.com  # get attack surface notifications 
exposed read port example.com           # read data from a feed. feed options: port, cve, login, subdomain
```

Simple to code against:

```go
package main

import (
  "github.com/hyperquack/exposed/sdk"
  "fmt"
)

func main(){
  client, err := sdk.Authenticate()
  if err != nil {
      panic(err)
  }

  resp, err := client.GetTargets()
  for _, target := resp.Hits {
      fmt.Println(target)
  }
}
```
