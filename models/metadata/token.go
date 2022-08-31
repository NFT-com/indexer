package metadata

type Token struct {
	Image           string      `json:"image"`
	ImageData       string      `json:"image_data"`
	ExternalURL     string      `json:"external_url"`
	Description     string      `json:"description"`
	Name            string      `json:"name"`
	BackgroundColor string      `json:"background_color"`
	AnimationURL    string      `json:"animation_url"`
	YoutubeURL      string      `json:"youtube_url"`
	Attributes      []Attribute `json:"attributes"`
}
