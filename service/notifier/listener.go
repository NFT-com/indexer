package notifier

type Listener interface {
	Notify(height uint64)
}
