package bootstrapper

type Config struct {
	ChainURL     string
	ChainType    string
	StandardType string
	Contract     string
	EventType    string
	StartIndex   int64
	EndIndex     int64
}
