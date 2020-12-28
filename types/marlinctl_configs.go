package types

type Project struct {
	Subscription  []string
	Version       string
	Storage       string
	Runtime       string
	ForcedRuntime bool
}

type Registry struct {
	Link    string
	Branch  string
	Local   string
	Enabled bool
}
