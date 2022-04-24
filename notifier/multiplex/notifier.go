package multiplex

import (
	"github.com/NFT-com/indexer/notifier"
)

type Notifier struct {
	listens []notifier.Listener
}

func NewNotifier(listens ...notifier.Listener) *Notifier {

	n := Notifier{
		listens: listens,
	}

	return &n
}

func (n *Notifier) Notify(height uint64) {
	for _, listen := range n.listens {
		listen.Notify(height)
	}
}
