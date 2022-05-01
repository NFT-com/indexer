package client

var DefaultConfig = Config{
	APIURL: "http://127.0.0.1:8080",
}

type Config struct {
	APIURL string
}

type Option func(*Config)

func WithAPIURL(url string) Option {
	return func(c *Config) {
		c.APIURL = url
	}
}
