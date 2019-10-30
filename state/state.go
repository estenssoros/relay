package state

var (
	Success        State = "success"
	Running        State = "running"
	Failed         State = "failed"
	Skipped        State = "skipped"
	Rescheduled    State = "rescheduleud"
	Retry          State = "retry"
	Queued         State = "queued"
	Pending        State = "pending"
	UpstreamFailed State = "upstream-failed"
	None           State = "none"
)

type State string
