# Exposed ("The Client")

The Exposed client is a command-line interface to the API.

> Valid keychain file is required to get started

It's simple to install:

```bash
go install github.com/privateducky/exposed-client@latest
```

And simple to use:

```bash
exposed targets                         # list domains you are monitoring
exposed start example.com               # start monitoring a new target
exposed stop example.com                # stop monitoring a target
exposed notify https://hooks.slack.com  # get push notifications 
exposed push port example.com 80        # add data to a feed
exposed pull port example.com           # pull data off a feed
```