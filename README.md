# Messenger

Messenger sends messages to a Discord webhook address while complying with the dynamic rate limiting (**User needs to manage rate limit imposed on channel/server, the library cannot detect if other webhooks also post messages to the same channel/server**).

## Usage

### Installing

```bash
go get github.com/qiyihuang/messenger
```

### Sending messages

```go
package main

import (
    "net/http"

    "github.com/qiyihuang/messenger"
)


func main() {
    hc := &http.Client{
        // Client config
    }
    url := "https://discord.com/api/webhooks/..."
    client, err := messenger.NewClient(hc, url)
    if err != nil {
        // handle when creating new client failed.
    }

    msgs := []messenger.Message{
        {
            Username: "Webhook username",
            Content:  "Message 1",
            Embeds: []messenger.Embed{
                // More fields please check go doc or Discord API.
                {Title: "Embed 1", Description: "Embed description 1"},
            },
        },
        {
            Username: "Webhook username",
            Content:  "Message 2",
        },
    }

    resp, err := client.Send(msgs)
    if err != nil {
        // handle when sending failed.
    }

    // ...
}
```

### Discord message limits

[Constants](https://pkg.go.dev/github.com/qiyihuang/messenger#pkg-constants) provided for managing message limits.
