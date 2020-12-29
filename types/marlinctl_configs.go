package types

type Project struct {
	Subscription  []string
	Version       string
	Storage       string
	Runtime       string
	ForcedRuntime bool
}

type Registry struct {
	Name    string
	Link    string
	Branch  string
	Local   string
	Enabled bool
}

type ReleaseJSON struct {
	JSONVersion int         `json:"json_version"`
	Data        interface{} `json:"data"`
}
