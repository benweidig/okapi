package okapi

// Info holds all data for the logging channel
type Info struct {
	Message   string
	Execution *ChangesetExecution
	Error     error
}

var notifier chan Info

// Notifier sets the channel for notifications
func Notifier(ch chan Info) {
	notifier = ch
}

func notify(msg string, ex *ChangesetExecution, err error) {
	if notifier == nil {
		return
	}

	notifier <- Info{
		Message:   msg,
		Execution: ex,
		Error:     err,
	}
}
