package messenger

// Author represents author object in an embed object.
type Author struct {
	Name    string `json:"name"`
	URL     string `json:"url"`
	IconURL string `json:"icon_url"`
}

// Field represents field object in an embed object.
type Field struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline,omitempty"`
}

// Footer represents footer object in an embed object.
type Footer struct {
	Text         string `json:"text"`
	IconURL      string `json:"icon_url,omitempty"`
	ProxyIconURL string `json:"proxy_icon_url,omitempty"`
}

// Embed represents an embed object in message object.
type Embed struct {
	Author      Author  `json:"author,omitempty"`
	Color       int     `json:"color,omitempty"`
	Title       string  `json:"title,omitempty"`
	Description string  `json:"description,omitempty"`
	URL         string  `json:"url,omitempty"`
	Fields      []Field `json:"fields,omitempty"`
	Footer      Footer  `json:"footer,omitempty"`
}

// Message represents a webhook message.
type Message struct {
	Username string  `json:"username,omitempty"`
	Embeds   []Embed `json:"embeds,omitempty"`
	Content  string  `json:"content,omitempty"`
}
