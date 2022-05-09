package notifier

type MultiNotifier struct {
	listens []Listener
}

func NewMultiNotifier(listens ...Listener) *MultiNotifier {

	m := MultiNotifier{
		listens: listens,
	}

	return &m
}

func (m *MultiNotifier) Notify(height uint64) {
	for _, listen := range m.listens {
		go listen.Notify(height)
	}
}
