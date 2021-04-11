# Messenger

Messenger is a library for sending Discord webhook. It sends messages to the a webhook address while complying with Discord's dynamic rate limiting (Channel and global limit is user's responsibility to manage since they cannot be detected by this library).

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

import "github.com/qiyihuang/messenger"


func main() {
    request := messenger.Request{
        Messages: []messenger.Message{
            {
                Username: "Webhook username",
                Content:  "Message 1",
                Embeds: []messenger.Embed{
                    // More fields please check exposed struct and Discord API
                    {Title: "Embed 1", Description: "Embed description 1"},
                },
            },
            {
                Username: "Webhook username",
                Content: "Message 2",
            },
        },
        URL: "https://discord.com/api/webhooks/...",
    }

    responses, err := request.Send()
    if err != nil {
        // Handle error.
    }

    // Continue...
}
```

### Discord message limits

Use [provided constants](https://pkg.go.dev/github.com/qiyihuang/messenger#pkg-constants) when you need to manage message limits.
