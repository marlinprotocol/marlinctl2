package runner

type Runner interface {
	PreRunSanity() error
	Download() error
	Prepare() error
	Create(runtimeArgs map[string]string) error
	Restart() error
	Recreate() error
	Destroy() error
	PostRun() error
	Status() error
	Logs(lines int) error
}
