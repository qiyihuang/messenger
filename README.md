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

[Messenger examples](https://github.com/qiyihuang/messenger/examples)