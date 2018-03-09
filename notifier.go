package okapi

type Info struct {
	Message string
	Record  *ChangesetExecution
	Error   error
}

var notifier chan Info

// Notifier sets the channel for notifications
func Notifier(ch chan Info) {
	notifier = ch
}

func notify(msg string, r *ChangesetExecution, err error) {
	if notifier == nil {
		return
	}

	notifier <- Info{
		Message: msg,
		Record:  r,
		Error:   err,
	}
}
