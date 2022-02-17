# Messenger

Messenger sends messages to a Discord webhook address while complying with the dynamic rate limiting.

**There is no way the libray can manage rate limit imposed on channel and server, user need to manage these limit in these cases (e.g. other webhooks also post messages to the same channel).**

## Usage

### In go.mod

```bash
require github.com/qiyihuang/messenger v0.4
```

### Send a message

```go
package main

import (
    "net/http"

    "github.com/qiyihuang/messenger"
)


func main() {
    client := &http.Client{
        // Client config
    }
    url := "https://discord.com/api/webhooks/..."
    msgs :=
        []messenger.Message{
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

    req, err := messenger.NewRequest(client, url, msgs)
    if err != nil {
        // handle when the request is invalid...
    }

    resp, err := req.Send()
    if err != nil {
        // ...
    }

    // ...
}
```

### Discord message limits

Use [these constants](https://pkg.go.dev/github.com/qiyihuang/messenger#pkg-constants) to manage message limits.
