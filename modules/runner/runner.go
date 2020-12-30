package runner

type Runner interface {
	PreRunSanity() error
	Download() error
	Prepare() error
	Create(runtimeArgs map[string]string) error
	Destroy() error
	PostRun() error
	Status() error
	Logs() error
}
