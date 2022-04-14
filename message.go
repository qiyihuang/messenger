package messenger

import "io"

// Message represents a webhook message.
type Message struct {
	Embeds   []Embed `json:"embeds,omitempty"`
	Files    []*File `json:"-"`
	Content  string  `json:"content,omitempty"`
	Username string  `json:"username,omitempty"`
}

// Embed represents an embed object in message object.
type Embed struct {
	Fields      []Field   `json:"fields,omitempty"`
	Author      Author    `json:"author,omitempty"`
	Footer      Footer    `json:"footer,omitempty"`
	Video       Video     `json:"video,omitempty"`
	Thumbnail   Thumbnail `json:"thumbnail,omitempty"`
	Image       Image     `json:"image,omitempty"`
	Title       string    `json:"title,omitempty"`
	Description string    `json:"description,omitempty"`
	URL         string    `json:"url,omitempty"`
	Timestamp   Timestamp `json:"timestamp,omitempty"`
	Color       int       `json:"color,omitempty"`
}

type File struct {
	// Name must include file extension (e.g. .jpg)
	Name string
	// Usually not required for common file types, check Discord docs if not working.
	// https://discord.com/developers/docs/reference#image-data/
	ContentType string
	Reader      io.Reader
}

// Timestamp represents the timestamp string in an embed object
// Format timestamp using .UTC().Format("2006-01-02T15:04:05-0700"),
// Discord will convert it to local time on display.
type Timestamp string

// Footer represents footer object in an embed object.
type Footer struct {
	Text         string `json:"text"`
	IconURL      string `json:"icon_url,omitempty"`
	ProxyIconURL string `json:"proxy_icon_url,omitempty"`
}

// Image represents the image object in an embed object.
type Image struct {
	URL      string `json:"url,omitempty"`
	ProxyURL string `json:"proxy_url,omitempty"`
	Height   int    `json:"height,omitempty"`
	Width    int    `json:"width,omitempty"`
}

// Thumbnail represents the thumbnail object in an embed object.
type Thumbnail struct {
	URL      string `json:"url,omitempty"`
	ProxyURL string `json:"proxy_url,omitempty"`
	Height   int    `json:"height,omitempty"`
	Width    int    `json:"width,omitempty"`
}

// Video represents the video object in an embed object.
type Video struct {
	URL      string `json:"url,omitempty"`
	ProxyURL string `json:"proxy_url,omitempty"`
	Height   int    `json:"height,omitempty"`
	Width    int    `json:"width,omitempty"`
}

// Author represents author object in an embed object.
type Author struct {
	Name         string `json:"name,omitempty"`
	URL          string `json:"url,omitempty"` // URL on the Author name field.
	IconURL      string `json:"icon_url,omitempty"`
	ProxyIconURL string `json:"proxy_icon_url,omitempty"`
}

// Field represents field object in an embed object.
type Field struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline,omitempty"`
}
