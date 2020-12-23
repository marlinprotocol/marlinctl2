package upstream

// Upstream configuration
type UpstreamConfig struct {
	RemoteVCS string
	HomeClone string
}

// Download upstream VCS to Home Clone location
func (c *UpstreamConfig) FetchUpstreamRegistry() (string, error) {
	return "", nil
}
