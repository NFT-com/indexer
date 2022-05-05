package inputs

type Addition struct {
	NodeURL      string `json:"node_url"`
	EventType    string `json:"event_type"`
	CollectionID string `json:"collection_id"`
	Owner        string `json:"owner"`
}
