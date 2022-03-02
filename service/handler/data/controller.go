package data

type Handler struct {
	store Store
}

func NewHandler(store Store) *Handler {
	c := Handler{
		store: store,
	}

	return &c
}
