package results

type Error struct {
	Message string `json:"errorMessage"`
	Type    string `json:"errorType"`
}
