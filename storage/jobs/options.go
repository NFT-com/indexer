package jobs

type updateOptions struct {
	statusMessage *string
}

type UpdateStatusOption func(*updateOptions)

func StatusMessage(statusMessage string) UpdateStatusOption {
	return func(options *updateOptions) {
		options.statusMessage = &statusMessage
	}
}
