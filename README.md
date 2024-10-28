# Exposed ("Client")

Exposed is a Continuous Threat Exposure utility that simplifies managing your attack surface.

## Usage

The client is a simple command-line interface (*and SDK!*) to create fully autonomous security tooling.

> A valid keychain file is required to get started

Simple to install:

```go
go install github.com/privateducky/exposed-client@latest
```

Simple to use:

```zsh
exposed targets                         # list the domains you are monitoring
exposed start example.com               # start monitoring a new target
exposed stop example.com                # stop monitoring a target
exposed notify https://hooks.slack.com  # get push notifications 
exposed push port example.com 80        # add data to a feed
exposed pull port example.com           # pull data off a feed
```