package results

type Error struct {
	Message string `json:"errorMessage"`
	Type    string `json:"errorType"`
}

func (e Error) Error() string {
	return e.Message
}
