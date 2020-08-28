package core

// Status types.
const (
	StatusSkipped  = "skipped"
	StatusBlocked  = "blocked"
	StatusDeclined = "declined"
	StatusWaiting  = "waiting_on_dependencies"
	StatusPending  = "pending"
	StatusRunning  = "running"
	StatusPassing  = "success"
	StatusFailing  = "failure"
	StatusKilled   = "killed"
	StatusError    = "error"
)
