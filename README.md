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

Use constants provided when you need to manage message limits.

```go
const (
    MessageEmbedNumLimit  = 10
    MessageContentLimit   = 2000
    EmbedTotalLimit       = 6000
    EmbedTitleLimit       = 256
    EmbedDescriptionLimit = 2048
    EmbedFieldNumLimit    = 25
    AuthorNameLimit       = 256
    FieldNameLimit        = 256
    FieldValueLimit       = 1024
    FooterTextLimit       = 2048
)
```
