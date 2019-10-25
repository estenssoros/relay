package state

var (
	Pending        State = "pending"
	Queued         State = "queued"
	Success        State = "success"
	Failed         State = "failed"
	UpstreamFailed State = "upstream-failed"
	Running        State = "running"
)

type State string
