package local

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/NFT-com/indexer/event"
)

type Dispatcher struct {
	url string
}

func New(url string) *Dispatcher {
	d := Dispatcher{
		url: url,
	}

	return &d
}

func (d *Dispatcher) Dispatch(function string, event *event.Event) error {
	payload, err := json.Marshal(event)
	if err != nil {
		// LOG
		return err
	}

	uri := fmt.Sprintf(d.url, function)
	resp, err := http.Post(uri, "application/json", bytes.NewReader(payload))
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusAccepted {
		return errors.New("failed to send request")
	}

	return nil
}
