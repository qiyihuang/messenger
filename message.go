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

// divideMessages breaks message into multiple messages depending on embed total
// character count and number of embeds.
func divideMessages(messages []Message) (msgs []Message) {
	for _, msg := range messages {
		dividedEmbeds := divideEmbeds(msg)
		// Create message for every embed chunk.
		for i, embeds := range dividedEmbeds {
			// First message contains content from original message.
			var content string
			var files []*File
			if i == 0 {
				content = msg.Content
				files = msg.Files
			}
			msgs = append(msgs, Message{Username: msg.Username, Embeds: embeds, Content: content, Files: files})
		}
	}
	return msgs
}

func divideEmbeds(msg Message) (dividedEmbeds [][]Embed) {
	var total int
	var startIndex int
	for i, e := range msg.Embeds {
		count := countEmbed(e)
		total += count
		// i will be included in next chunk.
		if total > EmbedTotalLimit || i-startIndex == MessageEmbedNumLimit {
			dividedEmbeds = append(dividedEmbeds, msg.Embeds[startIndex:i])
			startIndex = i
			// Set current count to initial total of next message.
			total = count
		}
	}
	// Add the last chunk.
	dividedEmbeds = append(dividedEmbeds, msg.Embeds[startIndex:])
	return
}

func countEmbed(e Embed) int {
	total := len(e.Title) + len(e.Description) + len(e.Author.Name) + len(e.Footer.Text)
	for _, field := range e.Fields {
		total += len(field.Name)
		total += len(field.Value)
	}
	return total
}
