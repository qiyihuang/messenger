# Messenger

Messenger is a library for sending Discord webhook. It can send single or multiple (WIP) messages to the same webhook address while complying with Discord's dynamic rate limiting (Channel and global limit cannot be detected by library so it's the user's responsibility to manage).

## Getting started

### Installing

```bash
go install github.com/qiyihuang/messenger
```

### Usage

Import tha package into your project.

```go
import "github.com/qiyihuang/messenger"
```

## Examples

### Send a message

```go
package main

import (
    "fmt"

    "github.com/qiyihuang/messenger"
)

func main() {
    request := messenger.Request{
        Msg: messenger.Message{
            Username: "Webhook username",
            Content:  "Hi content.",
            Embeds: []messenger.Embed{
                // More fields please check exposed struct and Discord API
                {Title: "Embed 1", Description: "Embed description 1"},
                {Title: "Embed 2", Description: "Embed description 2"},
            },
        },
        URL: "https://discord.com/api/webhooks/...",
    }

    resp, err := messenger.Send(request)
    if err != nil {
        // Handle error.
        fmt.Print(err)
    }

    // Continue...
}
```
