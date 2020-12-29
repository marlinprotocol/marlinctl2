package runner

type Runner interface {
	PreRunSanity() error
	Download() error
	Prepare() error
	Create() error
	Destroy() error
	PostRun() error
	Status() error
	Logs(tailLogs bool, prevLines int) error
}
